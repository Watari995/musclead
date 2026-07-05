package notificationusecase

// RegisterDeviceToken: アプリ起動時にクライアントから送られてくる FCM registration token を
// DeviceTokenRepository.Save (upsert) する usecase。
// NotificationHandler に POST /device-tokens のようなエンドポイントを追加して呼び出す想定。
