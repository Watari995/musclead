import type { Metadata, Viewport } from "next";
import {
  Architects_Daughter,
  Caveat,
  Geist_Mono,
  Yomogi,
} from "next/font/google";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";
import "./globals.css";
import { Providers } from "./providers";
import { AuthBootstrap } from "@/shared/auth/AuthBootstrap";
import { Header } from "@/features/user/ui/Header";
import { ThemePreferenceSync } from "@/features/user/ui/ThemePreferenceSync";

// 手描き(Excalidraw 風)デザインの UI フォント。Caveat = 見出し、
// Architects Daughter = 本文/UI(和文は Yomogi にフォールバック)。
const caveat = Caveat({
  variable: "--font-caveat",
  weight: "700",
  subsets: ["latin"],
});

const architectsDaughter = Architects_Daughter({
  variable: "--font-architects-daughter",
  weight: "400",
  subsets: ["latin"],
});

const yomogi = Yomogi({
  variable: "--font-yomogi",
  weight: "400",
  subsets: ["latin"],
});

const geistMono = Geist_Mono({
  variable: "--font-geist-mono",
  subsets: ["latin"],
});

export const metadata: Metadata = {
  title: "musclead",
  description: "筋トレ・食事・体重 一元管理 SaaS",
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  viewportFit: "cover",
};

export default async function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  const locale = await getLocale();
  const messages = await getMessages();

  return (
    <html
      lang={locale}
      className={`${caveat.variable} ${architectsDaughter.variable} ${yomogi.variable} ${geistMono.variable} h-full antialiased`}
      suppressHydrationWarning
    >
      <body className="min-h-full flex flex-col overflow-x-hidden">
        {/* .icon-sketchy / 手描き輪郭の乱数ゆらぎに使う SVG filter。ページ内に 1 度だけ置く。 */}
        <svg width="0" height="0" style={{ position: "absolute" }} aria-hidden>
          <filter id="sketchy" x="-30%" y="-30%" width="160%" height="160%">
            <feTurbulence
              type="fractalNoise"
              baseFrequency="0.045"
              numOctaves={2}
              seed={7}
              result="n"
            />
            <feDisplacementMap in="SourceGraphic" in2="n" scale={2.4} />
          </filter>
        </svg>
        <NextIntlClientProvider locale={locale} messages={messages}>
        <Providers>
          <AuthBootstrap />
          <ThemePreferenceSync />
          <Header />
          <main className="flex-1 w-full max-w-5xl mx-auto px-4 sm:px-6 py-6 sm:py-8">
            {children}
          </main>
          <footer className="border-t border-[var(--color-line)] mt-12">
            <div className="w-full max-w-5xl mx-auto px-4 sm:px-6 py-6 text-xs text-[var(--color-ink-muted)] flex justify-between">
              <span>© musclead</span>
              <span>Beta</span>
            </div>
          </footer>
        </Providers>
        </NextIntlClientProvider>
      </body>
    </html>
  );
}
