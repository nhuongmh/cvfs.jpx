"use client";
import ProposalCheck from "./review";
import { Title } from "components/title";
import React from "react";
import { useSearchParams, usePathname } from 'next/navigation'

export default function Process() {
  const searchParams = useSearchParams();
  const pathname = usePathname();
  if (pathname === null || searchParams === null) {
    return null
  }
  const groupId = pathname.split('/').pop() || searchParams.get('id') || "";
  return (
    <div className="container px-8 mx-auto mt-16 lg:mt-32 ">
      <Title>Process Proposal Cards</Title>
      <ProposalCheck group={groupId as string || ""} />
    </div>
  );
}
