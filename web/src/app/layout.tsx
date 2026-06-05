import type { Metadata } from "next";
import { Geist, Geist_Mono } from "next/font/google";
import "./globals.css";
import { Providers } from "./providers";
import { AuthBootstrap } from "@/shared/auth/AuthBootstrap";
import { Header } from "@/features/user/ui/Header";

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

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html
      lang="ja"
      className={`${geistSans.variable} ${geistMono.variable} h-full antialiased`}
    >
      <body className="min-h-full flex flex-col">
        <Providers>
          <AuthBootstrap />
          <Header />
          <main className="flex-1 w-full max-w-5xl mx-auto px-5 sm:px-6 py-8">
            {children}
          </main>
          <footer className="border-t border-[var(--color-line)] mt-12">
            <div className="w-full max-w-5xl mx-auto px-5 sm:px-6 py-6 text-xs text-[var(--color-ink-muted)] flex justify-between">
              <span>© musclead</span>
              <span>Beta</span>
            </div>
          </footer>
        </Providers>
      </body>
    </html>
  );
}
