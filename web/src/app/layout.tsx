import type { Metadata, Viewport } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import { NextIntlClientProvider } from "next-intl";
import { getLocale, getMessages } from "next-intl/server";
import "./globals.css";
import { Providers } from "./providers";
import { AuthBootstrap } from "@/shared/auth/AuthBootstrap";
import { Header } from "@/features/user/ui/Header";
import { ThemePreferenceSync } from "@/features/user/ui/ThemePreferenceSync";

const geistSans = Geist({
  variable: "--font-geist-sans",
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
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
      suppressHydrationWarning
    >
      <body className="min-h-full flex flex-col overflow-x-hidden">
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
