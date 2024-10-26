"use client";
import CardLearn from "./learn";
import { Title } from "components/title";
import React from "react";

export default function Process() {
  return (
    <div className="container px-8 mx-auto mt-16 lg:mt-32 ">
      <Title>Learn Cards</Title>
      <CardLearn />
    </div>
  );
}
