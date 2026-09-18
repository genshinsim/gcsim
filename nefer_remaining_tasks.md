# Nefer: оставшиеся задачи

Этот файл составлен по текущим состояниям из:

- [Nefer PR checklist.md](Nefer%20PR%20checklist.md)
- [nefer_implementation_plan.md](nefer_implementation_plan.md)
- [nefer_implementation_progress.md](nefer_implementation_progress.md)
- [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md)
- [nefer_readiness_assessment.md](nefer_readiness_assessment.md)
- [nefer_frames_google_sheets.md](nefer_frames_google_sheets.md)
- [nefer_ingame_observations.md](nefer_ingame_observations.md)
- [nefer_lunaris_10000122.md](nefer_lunaris_10000122.md)

Статус на момент составления: базовый игровой цикл и констелляции реализованы, но реализация ещё не готова считаться полностью combat-accurate.

## 1. Приоритетные задачи реализации

### 1.1. Закрыть seed- и absorb-механику

- [ ] Определить и реализовать точный радиус поглощения Seeds of Deceit вместо текущего `seedAbsorbRadius = 6`.
- [ ] Проверить правило выбора поглощаемых seed: учитывать ли вертикальное расстояние, направление движения Slither и другие ограничения.
- [ ] Проверить порядок удаления при общем лимите из 5 Dendro Core/Seed и решить, достаточно ли текущего общего `GadgetTypDendroCore`.
- [ ] Определить, должны ли seed иметь отдельную collision/interaction-логику, а не только быть поглощаемыми гаджетами.
- [ ] Повторно проверить lifetime seed после конвертации, сохранение seed после окончания P1 и поведение общего лимита.
- [ ] Проверить источник и точность порога выносливости для входа в Slither; текущая модель разделяет readiness и фактическое потребление, но численное значение не подтверждено источником.

### 1.2. Закрыть timing и geometry Phantasm/C6

- [ ] Проверить точные кадры четвёртого и пятого ударов Phantasm.
- [ ] Определить точный момент дополнительного C6-удара после завершения Phantasm.
- [ ] Подтвердить тайминги и поведение цепочек `Phantasm -> Skill/Burst/Dash/Jump/Walk/Swap`.
- [ ] Подтвердить AoE и hitbox для всех ударов Phantasm.
- [ ] Подтвердить hitbox дополнительного C6-удара.
- [ ] Провести отдельную интерпретацию страницы Phantasm CA из Workbook 2 и использовать её только после подтверждения смысла измерений.

### 1.3. Добавить явные PoiseDMG и Hitlag

- [ ] Найти source-backed значения и добавить `PoiseDMG`, `HitlagFactor` и `HitlagHaltFrames` для Normal 1-4.
- [ ] Добавить те же значения для Low Plunge и High Plunge.
- [ ] Добавить значения для обычной Charged Attack.
- [ ] Добавить значения для Skill.
- [ ] Добавить значения для всех пяти обычных ударов Phantasm.
- [ ] Добавить значения для C6 End Hit.
- [ ] Добавить значения для обоих Burst-ударов.

## 2. Уточнение кадров Normal Attack и Charged Attack

### 2.1. Normal Attack

- [ ] Уточнить смешанные hitmark-диапазоны N1 `10-13`, N2 `8-11`, N4 `22-24`.
- [ ] Уточнить переходы N1 -> N2 `15`, N1 -> CA `40-41`, N1 -> Walk `32-33`.
- [ ] Уточнить переходы N2 -> N3 `18-21` и N2 -> Phantasm CA `38-39`.
- [ ] Уточнить переходы N3 -> N4 `48-53` и N3 -> Normal CA `61-62`.
- [ ] Уточнить переходы N4 -> N1 `50-54` и N4 -> Walk `62-64`.
- [ ] Получить или явно задокументировать отсутствующие строки переходов Normal Attack -> Skill/Burst/Dash/Jump/Swap.
- [ ] Определить отдельные hitbox для N1, N2, обоих попаданий N3, N4 и plunge вместо текущих приближённых геометрий.

### 2.2. Обычная Charged Attack и Slither

- [ ] Проверить смешанный переход CA -> CA `48-49`.
- [ ] Найти данные для CA -> Skill/Burst/Dash/Jump.
- [ ] Уточнить CA -> Swap и принятую модель полного no-hold Slither-release маршрута.
- [ ] Проверить, действительно ли embedded Slither-entry имеет `0f`, либо нужен отдельный внутренний timing.
- [ ] Уточнить минимальный кадр отмены Slither `24`.
- [ ] Подтвердить cadence движения Slither `1f`.
- [ ] Подтвердить расстояние движения `0.1` за тик.
- [ ] Подтвердить cadence расхода выносливости `1f`.
- [ ] Подтвердить hitbox выхода из обычной Charged Attack.

## 3. Уточнение Skill и Burst

### 3.1. Skill

- [ ] Уточнить переход Skill -> Normal CA `47-49`.
- [ ] Уточнить переход Skill -> Phantasm CA `28-29`.
- [ ] Уточнить переход Skill -> Skill `69-74`.
- [ ] Уточнить переходы Skill -> Burst `26-28`, Skill -> Dash `37-39`, Skill -> Swap `24-25`.
- [ ] Закрыть точное соответствие datamine lock-shape `CircleLockEnemyR8H6HC`/`CircleLockEnemyR15H10HC` фактическому hitbox Skill.
- [ ] Проверить, привязано ли создание частиц к hit callback на кадре 24 во всех ситуациях.
- [ ] Проверить particle ICD `0.2s` на нескольких целях и повторных попаданиях.
- [ ] Проверить распределение частиц: 66% для 3 частиц и 33% для 2 частиц.

### 3.2. Burst

- [ ] Уточнить hitmark первого Burst-удара `99-101` вместо текущего значения.
- [ ] Уточнить hitmark второго Burst-удара `40-44` вместо текущего значения.
- [ ] Проверить точный момент списания 60 энергии; текущая немедленная модель не подтверждена смешанными наблюдениями `6-7`.
- [ ] Уточнить переходы Burst -> N1 `118-119`, Burst -> CA `129-131`, Burst -> Phantasm CA `123-124`.
- [ ] Уточнить переходы Burst -> Dash `119-120`, Burst -> Jump `119-120`, Burst -> Walk `118-119`, Burst -> Swap `117-118`.
- [ ] Закрыть точные hitbox обоих Burst-ударов по datamine lock-shape.

## 4. Геометрия и area-эффекты

- [ ] Получить source-backed geometry для всех Normal Attack и plunge hitbox.
- [ ] Получить source-backed geometry для обычной Charged Attack.
- [ ] Получить source-backed AoE radius для каждого Phantasm hit.
- [ ] Получить source-backed AoE для C6 End Hit.
- [ ] Определить точную область C4 nearby-opponent RES shred вместо текущего `c4NearbyRadius = 10`.
- [ ] Проверить, соответствует ли текущая circle-based реализация реальным lock shapes и вертикальным ограничениям.

## 5. Финальная валидация уже реализованных назначений

- [ ] Провести package-wide проверку ICD для Normal Attack, plunge, Skill, Burst, обычной CA, Phantasm и Lunar-Bloom-tagged ударов.
- [ ] Провести package-wide проверку StrikeType, включая blunt Skill hit.
- [ ] Проверить назначенные durability values на всех attack paths.
- [ ] Проверить particle behavior после закрытия Skill timing/edge cases.
- [ ] Проверить C1/C2 порядок применения множителей Phantasm и убедиться, что поздние additive reaction terms не попадают под Veil multiplier.
- [ ] Проверить C6: 85% EM второй стадии, дополнительный 120% EM удар и 15% Lunar-Bloom elevation.
- [ ] Проверить одноразовый Veil snapshot перед первым ударом Phantasm и его применение ко всей последовательности.
- [ ] Проверить Veil stack duration, C2 cap в 5 стеков и замену самого старого стека при заполнении.
- [ ] Повторить smoke/regression-проверки для:
  - [ ] Phantasm только с Moonridge Dew.
  - [ ] Конвертации старых Dendro Core в Seed и получения Veil на первой CA.
  - [ ] Отмены Phantasm Burst/Dash после первого удара.
  - [ ] Сброса Shadow Dance и Phantasm charges после swap.
  - [ ] P1/P2 gating и обновления Verdant Dew.

## 6. Интеграция, генерация и документация

- [ ] После изменений запустить генерацию и проверить синхронизацию:
  - [ ] `internal/characters/nefer/nefer_gen.go`
  - [ ] `pkg/core/keys/keys_char_gen.go`
  - [ ] `pkg/simulation/imports_char_gen.go`
  - [ ] `ui/packages/db/src/Data/char_data.generated.json`
  - [ ] `ui/packages/ui/src/Data/char_data.generated.json`
  - [ ] generated docs для Nefer
- [ ] Проверить, что [ui/packages/docs/docs/reference/characters/nefer.md](ui/packages/docs/docs/reference/characters/nefer.md) соответствует текущей реализации.
- [ ] Обновлять [nefer_implementation_progress.md](nefer_implementation_progress.md) только для завершённых задач.
- [ ] Поддерживать [nefer_inexact_implementation_register.md](nefer_inexact_implementation_register.md) как основной список оставшихся технических gaps.
- [ ] Поддерживать [nefer_ingame_observations.md](nefer_ingame_observations.md) только как журнал наблюдений, не добавляя туда implementation policy.
- [ ] Обновить [Nefer PR checklist.md](Nefer%20PR%20checklist.md) после закрытия каждого крупного workstream.
- [ ] Выполнить compile/test-проверку после изменений и зафиксировать результат в progress log.

## 7. Необязательные задачи

- [ ] Уточнить Nefer-specific Yelan N0 hook: текущие `10f` остаются best-fit предположением.
- [ ] Уточнить Nefer-specific Xingqiu N0 hook: предполагаемые `11f` не привязаны к подтверждённому gcsim hook.
- [ ] Решить, требуется ли отдельная поддержка Xianyun Plunge.
- [ ] Если поддержка нужна, исследовать и добавить Nefer-specific Xianyun Plunge timing/interaction.
- [ ] Проверить страницу DashJump и определить, можно ли безопасно использовать её значения для gcsim hooks.

## Рекомендуемый порядок выполнения

1. Seed absorption, seed lifetime/cap и порог входа в Slither.
2. Phantasm и C6 timing/geometry.
3. Skill/Burst timing и geometry.
4. Явные PoiseDMG/Hitlag mappings.
5. Финальная проверка ICD, StrikeType, durability и particles.
6. Генерация артефактов, тесты и обновление PR-документации.

## Что уже не считается оставшейся задачей

- Регистрация персонажа, ключи, импорты, shortcuts и базовый data pipeline.
- Базовая реализация Normal Attack, Charged Attack, Skill, Burst и plunge.
- Shadow Dance, Slither loop, Phantasm charges и swap reset.
- Seed conversion window и сохранение seed после окончания P1.
- Veil stacks, C1, C2, C4 и C6 в текущей принятой интерпретации.
- Базовое назначение ICD, StrikeType, durability и частиц; для них остаётся только финальная валидация.
