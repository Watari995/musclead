<div align="center">

# 💪 musclead

**Training. Meals. Weight. One app.**

Stop juggling three apps to track one body.

</div>

---

## Why musclead?

You lift, you eat, you weigh in. Every serious trainee already does all three —
but the data lives in three different apps that don't talk to each other.
You can't see how last week's calorie deficit hit your bench numbers.
You can't tell whether the scale moved because of muscle or because of yesterday's ramen.

**musclead puts all three on one timeline, so the story of your body actually adds up.**

## What you can do

- 🏋️ **Log workouts** — sets, reps, weight. Build it once as a routine, reuse it forever.
- ⚖️ **Track your weight** — daily check-ins with body fat and muscle mass, trend over time.
- 🍱 **Log meals** — PFC macros tracked automatically, per meal and per day.
- 🔐 **Your account, your data** — personal, private, no social feed.

## Why you'll actually use it

- **All-in-one.** No more "log workout in Strong, weight in Health, food in MyFitnessPal."
- **Fast and minimal.** No ads, no nags, no streak shame.
- **Built for trainees who are sick of bloat.** If you've ever closed an app because it asked you to rate it mid-set, you're the target user.

## Platforms

| Platform | Status |
|---|---|
| **iOS (App Store)** | [配信中](https://apps.apple.com/jp/app/musclead/id6781308658) |
| **Web** | https://app.musclead.com |

<details>
<summary>👩‍💻 For developers</summary>

**Stack**

| Layer | Tech |
|---|---|
| Backend | Go (net/http REST + OpenAPI/swag) + gorp + DDD |
| Frontend (Web) | Next.js + React + Tailwind v4 |
| iOS | Flutter 3.x (Riverpod v2 / go_router) |
| DB | PostgreSQL |
| Hosting | Oracle VM (API) + Vercel (Web) |

See [`CLAUDE.md`](CLAUDE.md), [`docs/domain-model.md`](docs/domain-model.md), and the [ADRs](docs/adr/) for the full picture.

```bash
# Backend → :8080
make db-up && make migrate-up && make dev

# Web → :3000
cd web && npm install && npm run dev

# iOS (requires fvm)
cd mobile && fvm flutter pub get && fvm flutter run
```

</details>
