import 'package:cookie_jar/cookie_jar.dart';
import 'package:dio/dio.dart';
import 'package:dio_cookie_manager/dio_cookie_manager.dart';

import '../config/app_config.dart';

/// 認証 interceptor を除いた素の Dio を構築する。
/// （AuthInterceptor は自分自身（dio）を retry に使うため、provider 側で後付けする）
Dio buildDio({required CookieJar cookieJar}) {
  final dio = Dio(
    BaseOptions(
      baseUrl: AppConfig.current.apiBaseUrl,
      connectTimeout: const Duration(seconds: 15),
      receiveTimeout: const Duration(seconds: 20),
      headers: {'Accept': 'application/json'},
      contentType: Headers.jsonContentType,
    ),
  );
  dio.interceptors.add(CookieManager(cookieJar));
  return dio;
}
