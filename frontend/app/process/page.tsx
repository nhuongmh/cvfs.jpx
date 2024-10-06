"use client";
import { ArrowTopRightOnSquareIcon } from "@heroicons/react/20/solid";
import Link from "next/link";
import SpacedRepetitionApp from "./review";
import { Title } from "@components/title";
import React from "react";
const steps: {
  name: string;
  description: string | React.ReactNode;
  cta?: React.ReactNode;
}[] = [
  {
    name: "Create a new Redis database on Upstash",
    description: (
      <>
        Upstash offers a serverless Redis database with a generous free tier of up to 10,000 requests per day. That's
        more than enough.
        <br />
        Click the button below to sign up and create a new Redis database on Upstash.
      </>
    ),
    cta: (
      <Link
        href="https://console.upstash.com/redis"
        className="flex items-center justify-center w-full gap-2 px-4 py-2 text-sm text-center transition-all duration-150 rounded text-zinc-800 hover:text-zinc-100 bg-zinc-200 hover:bg-transparent ring-1 ring-zinc-100"
      >
        <span>Create Database</span>
        <ArrowTopRightOnSquareIcon className="w-4 h-4" />
      </Link>
    ),
  },
  {
    name: "Copy the REST connection credentials",
    description: (
      <p>
        After creating the database, scroll to the bottom and make a note of <code>UPSTASH_REDIS_REST_URL</code> and{" "}
        <code>UPSTASH_REDIS_REST_TOKEN</code>, you need them in the next step
      </p>
    ),
  },
  {
    name: "Deploy to Vercel",
    description: "Deploy the app to Vercel and paste the connection credentials into the environment variables.",
    cta: (
      <Link
        href="https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2Fchronark%2Fenvshare&env=UPSTASH_REDIS_REST_URL,UPSTASH_REDIS_REST_TOKEN&demo-title=Share%20Environment%20Variables%20Securely&demo-url=https%3A%2F%2Fcryptic.vercel.app"
        className="flex items-center justify-center w-full gap-2 px-4 py-2 text-sm text-center transition-all duration-150 rounded text-zinc-800 hover:text-zinc-100 bg-zinc-200 hover:bg-transparent ring-1 ring-zinc-100"
      >
        <span>Deploy</span>
        <ArrowTopRightOnSquareIcon className="w-4 h-4" />
      </Link>
    ),
  },
];

export default function Process() {
  return (
    <div className="container px-8 mx-auto mt-16 lg:mt-32 ">
      <Title>Deploy EnvShare for Free</Title>
      <SpacedRepetitionApp />
    </div>
  );
}
