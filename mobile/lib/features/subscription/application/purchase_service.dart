import 'dart:async';

import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:in_app_purchase/in_app_purchase.dart';

/// StoreKit (IAP) サービスの雛形 — **Phase 3**。
///
/// 前提（あなたのコンソール作業）:
/// - App Store Connect で自動更新サブスク商品 [proProductId] を登録
/// - Sandbox テスターで検証
/// - 購入完了後はサーバ側で App Store レシート(JWS)検証 +
///   App Store Server Notifications V2 を実装し subscriptions を更新する
///
/// MVP ではアプリ内に購入導線を出さない（App Store Review 3.1.1）。
/// 本サービスは Phase 3 で課金を有効化する際の足場。
class PurchaseService {
  PurchaseService(this._iap);

  final InAppPurchase _iap;

  /// App Store Connect に登録する商品 ID（要登録）。
  static const String proProductId = 'com.musclead.pro.monthly';

  StreamSubscription<List<PurchaseDetails>>? _sub;

  Future<bool> isAvailable() => _iap.isAvailable();

  /// 購入ストリームを購読する。purchased/restored はサーバ検証後に確定させる。
  void listen({required Future<void> Function(PurchaseDetails) onPurchase}) {
    _sub ??= _iap.purchaseStream.listen((purchases) async {
      for (final p in purchases) {
        if (p.status == PurchaseStatus.purchased ||
            p.status == PurchaseStatus.restored) {
          await onPurchase(p); // TODO(phase3): サーバでレシート検証
        }
        if (p.pendingCompletePurchase) {
          await _iap.completePurchase(p);
        }
      }
    });
  }

  Future<ProductDetails?> proProduct() async {
    final res = await _iap.queryProductDetails({proProductId});
    if (res.productDetails.isEmpty) return null;
    return res.productDetails.first;
  }

  Future<bool> buyPro(ProductDetails product) => _iap.buyNonConsumable(
    purchaseParam: PurchaseParam(productDetails: product),
  );

  void dispose() => _sub?.cancel();
}

final purchaseServiceProvider = Provider<PurchaseService>((ref) {
  final service = PurchaseService(InAppPurchase.instance);
  ref.onDispose(service.dispose);
  return service;
});
