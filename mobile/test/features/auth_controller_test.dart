import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:musclead/core/providers/core_providers.dart';
import 'package:musclead/features/auth/application/auth_controller.dart';
import 'package:musclead/features/auth/data/auth_repository.dart';

import '../support/fakes.dart';

class MockAuthRepository extends Mock implements AuthRepository {}

Future<void> _flushMicrotasks() => Future<void>.delayed(Duration.zero);

void main() {
  late MockAuthRepository repo;
  late FakeTokenStore store;

  ProviderContainer makeContainer() {
    final container = ProviderContainer(
      overrides: [
        authRepositoryProvider.overrideWithValue(repo),
        tokenStoreProvider.overrideWithValue(store),
      ],
    );
    addTearDown(container.dispose);
    return container;
  }

  setUp(() {
    repo = MockAuthRepository();
    store = FakeTokenStore();
  });

  test('トークン無しで起動 → unauthenticated', () async {
    final container = makeContainer();
    container.read(authControllerProvider); // build をトリガ
    await _flushMicrotasks();
    expect(container.read(authControllerProvider), AuthStatus.unauthenticated);
  });

  test('login 成功 → authenticated', () async {
    when(() => repo.login(any(), any())).thenAnswer((_) async {});
    final container = makeContainer();
    final notifier = container.read(authControllerProvider.notifier);
    await _flushMicrotasks();

    await notifier.login('a@example.com', 'pw');

    expect(container.read(authControllerProvider), AuthStatus.authenticated);
    verify(() => repo.login('a@example.com', 'pw')).called(1);
  });

  test('logout → unauthenticated', () async {
    when(() => repo.logout()).thenAnswer((_) async {});
    await store.writeAccessToken('token'); // 起動時は authenticated
    final container = makeContainer();
    final notifier = container.read(authControllerProvider.notifier);
    await _flushMicrotasks();
    expect(container.read(authControllerProvider), AuthStatus.authenticated);

    await notifier.logout();

    expect(container.read(authControllerProvider), AuthStatus.unauthenticated);
  });
}
