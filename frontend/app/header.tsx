"use client";
import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { siteConfig } from "@/config/site";

export const Header: React.FC = () => {
  const pathname = usePathname();
  return (
    <header className="top-0 z-30 w-full px-4 sm:fixed">
      <div className="container mx-auto">
        <div className="flex flex-col items-center justify-between gap-2 pt-6 sm:h-20 sm:flex-row sm:pt-0">
          <Link href="/" className="text-2xl font-semibold duration-150 text-zinc-100 hover:text-white">
            .langfi
          </Link>
          {/* Desktop navigation */}
          <nav className="flex items-center grow">
            <ul className="flex flex-wrap items-center justify-end gap-4 grow">
              {siteConfig.navItems.map((item) => (
                <li className="" key={item.href}>
                  <Link
                    className={`flex items-center px-3 py-2 duration-150 text-sm sm:text-base  hover:text-zinc-50
                    ${pathname === item.href ? "text-zinc-200" : "text-zinc-400"}`}
                    href={item.href}
                    target={item.external ? "_blank" : undefined}
                    rel={item.external ? "noopener noreferrer" : undefined}
                  >
                    {item.name}
                  </Link>
                </li>
              ))}
            </ul>
          </nav>
        </div>
      </div>

      {/* Fancy fading bottom border */}
    </header>
  );
};
