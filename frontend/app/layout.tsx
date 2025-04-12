"use client";
import "@/styles/globals.css";
import { Inter } from "next/font/google";
import Link from "next/link";
import { Analytics } from "@/components/analytics";
import { Header } from "./header";
import { Providers } from "./providers";

const inter = Inter({ subsets: ["latin"], variable: "--font-inter" });

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en" className={inter.variable}>
      <head />
      <body className="dark text-foreground bg-background relative min-h-screen from-zinc-900/50 to-zinc-700/30">
        {
          // Not everyone will want to host envshare on Vercel, so it makes sense to make this opt-in.
          process.env.ENABLE_VERCEL_ANALYTICS ? <Analytics /> : null
        }
        <Providers themeProps={{ attribute: "class", defaultTheme: "dark", children: children }}>
        <Header />

        <main className=" min-h-[80vh] ">{children}</main>

        <footer className="bottom-0 border-t inset-2x-0 border-zinc-500/10">
          <div className="flex flex-col gap-1 px-6 py-12 mx-auto text-xs text-center text-zinc-700 max-w-7xl lg:px-8">
            <p>
              Built by{" "}
              <Link href="https://twitter.com/nhuongmh" className="font-semibold duration-150 hover:text-zinc-200">
                @nhuongmh
              </Link>
            </p>
          </div>
        </footer>
        </Providers>
      </body>
    </html>
  );
}
