#!/usr/bin/env bash
# musclead 本番 ECS/ALB/ログ 読み取り診断。
# 読み取り専用の AWS 認証情報(aws configure 済み)が前提。書き込みは一切しない。
#
# 使い方:
#   bash docs/ops/diagnose-prod.sh
#   AWS_PROFILE=musclead-diag bash docs/ops/diagnose-prod.sh   # 名前付きプロファイル利用時
set -uo pipefail

export AWS_REGION="${AWS_REGION:-ap-northeast-1}"
CLUSTER="musclead-cluster"
SERVICE="musclead-server-service"
TG_NAME="musclead-server-tg"
LOG_GROUP="/musclead/ecs/server"

echo "===== 0. 認証確認 ====="
aws sts get-caller-identity || { echo "認証エラー: aws configure を確認"; exit 1; }

echo; echo "===== 1. ECS service: desired/running, deployment 設定, 直近イベント ====="
aws ecs describe-services --cluster "$CLUSTER" --services "$SERVICE" \
  --query 'services[0].{desired:desiredCount,running:runningCount,pending:pendingCount,deploymentConfiguration:deploymentConfiguration,deployments:deployments[].{status:status,rolloutState:rolloutState,desired:desiredCount,running:runningCount,failed:failedTasks},events:events[0:12].[createdAt,message]}' \
  --output json

echo; echo "===== 2. 稼働中タスク ====="
aws ecs list-tasks --cluster "$CLUSTER" --service-name "$SERVICE" --desired-status RUNNING --output text --query 'taskArns'

echo; echo "===== 3. 停止タスクと停止理由(Spot 回収/クラッシュ判定) ====="
STOPPED=$(aws ecs list-tasks --cluster "$CLUSTER" --service-name "$SERVICE" --desired-status STOPPED --query 'taskArns' --output text)
if [ -n "$STOPPED" ] && [ "$STOPPED" != "None" ]; then
  aws ecs describe-tasks --cluster "$CLUSTER" --tasks $STOPPED \
    --query 'tasks[].{stoppedAt:stoppedAt,stopCode:stopCode,stoppedReason:stoppedReason,containers:containers[].{name:name,exitCode:exitCode,reason:reason}}' \
    --output json
else
  echo "(停止タスクなし)"
fi

echo; echo "===== 4. ALB ターゲットの health 状態 ====="
TG_ARN=$(aws elbv2 describe-target-groups --names "$TG_NAME" --query 'TargetGroups[0].TargetGroupArn' --output text)
aws elbv2 describe-target-health --target-group-arn "$TG_ARN" \
  --query 'TargetHealthDescriptions[].{target:Target.Id,port:Target.Port,state:TargetHealth.State,reason:TargetHealth.Reason,description:TargetHealth.Description}' \
  --output json

echo; echo "===== 5. 直近1時間のログ(panic/error/SIGTERM/Spot 等) ====="
aws logs filter-log-events --log-group-name "$LOG_GROUP" \
  --start-time "$(( ($(date +%s) - 3600) * 1000 ))" \
  --filter-pattern '?panic ?fatal ?error ?ERROR ?SIGTERM ?interrupt ?reclaim' \
  --query 'events[-50:].[message]' --output text 2>/dev/null || echo "(該当ログなし、またはロググループ名/権限を確認)"

echo; echo "===== 完了 ====="
