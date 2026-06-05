#!/usr/bin/env node
/**
 * Backend の /swagger/doc.json (Swagger 2.0) を取得し、
 * OpenAPI 3 に変換 → X-User-ID ヘッダ params を除去 → TS 型を生成。
 *
 * X-User-ID は openapi-fetch のミドルウェアで自動付与するため、
 * 呼び出し側の型からは消す。
 */
import { writeFileSync } from "node:fs";
import { mkdtempSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { execSync } from "node:child_process";
import converter from "swagger2openapi";

const SWAGGER_URL =
  process.env.SWAGGER_URL ?? "http://localhost:8080/swagger/doc.json";
const OUT = "src/shared/api/schema.ts";

const res = await fetch(SWAGGER_URL);
if (!res.ok) throw new Error(`fetch ${SWAGGER_URL}: HTTP ${res.status}`);
const swagger2 = await res.json();

const { openapi: openapi3 } = await converter.convertObj(swagger2, {});

for (const path of Object.values(openapi3.paths ?? {})) {
  for (const method of Object.values(path)) {
    if (!method?.parameters) continue;
    method.parameters = method.parameters.filter(
      (p) => !(p.in === "header" && p.name === "X-User-ID"),
    );
  }
}

const tmp = mkdtempSync(join(tmpdir(), "openapi-"));
const specPath = join(tmp, "openapi3.json");
writeFileSync(specPath, JSON.stringify(openapi3));

execSync(`npx openapi-typescript ${specPath} -o ${OUT}`, { stdio: "inherit" });
