"use client";
import React from "react";
import Link from "next/link";
import { usePathname } from "next/navigation";
import { siteConfig } from "@/config/site";

export const Header: React.FC = () => {
  const pathname = usePathname();
  return (
    <header className="top-0 z-30 px-4 sm:fixed">
      <div className="container mx-0">
      <div className="flex items-center justify-between pt-6 sm:h-20 sm:pt-0">
        {/* Collapsible menu icon */}
        <div className="relative">
        <button
          className="text-2xl font-semibold duration-150 text-zinc-100 hover:text-white"
          onClick={() => {
          const menu = document.getElementById("mobile-menu");
          if (menu) {
            menu.classList.toggle("hidden");
          }
          }}
        >
          ☰
        </button>
        <nav
          id="mobile-menu"
          className="absolute left-0 top-full mt-2 hidden w-48 bg-zinc-800 rounded shadow-lg"
        >
          <ul className="flex flex-col gap-2 p-4">
          {siteConfig.navItems.map((item) => (
            <li key={item.href}>
            <Link
              className={`block px-3 py-2 duration-150 text-sm sm:text-base hover:text-zinc-50
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
      </div>
    </header>
  );
};
