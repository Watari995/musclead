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
- ⚖️ **Track your weight** — daily check-ins, trend over time, no drama.
- 🍱 **Log meals** — what you actually ate, not what an algorithm guessed.
- 🔐 **Your account, your data** — personal, private, no social feed.

## Why you'll actually use it

- **All-in-one.** No more "log workout in Strong, weight in Health, food in MyFitnessPal."
- **Fast and minimal.** No ads, no nags, no streak shame, no premium upsell.
- **Built for trainees who are sick of bloat.** If you've ever closed an app because it asked you to rate it mid-set, you're the target user.

<details>
<summary>👩‍💻 For developers</summary>

Go backend (DDD + modular monolith) + Next.js frontend, typed end-to-end via OpenAPI.
See [`CLAUDE.md`](CLAUDE.md), [`docs/domain-model.md`](docs/domain-model.md), and the [ADRs](docs/adr/) for the full picture.

```bash
make db-up && make migrate-up && make dev   # backend → :8080
cd web && npm install && npm run dev        # frontend → :3000
```

</details>
