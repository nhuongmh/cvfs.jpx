"use client";
import { Title } from "components/title";
import React from "react";
import { useRouter } from 'next/router'

export default function Process() {
    const router = useRouter()
    return (
        <div className="container px-8 mx-auto mt-16 lg:mt-32 ">
            <Title>Learn Cards for {router.query.id}</Title>
        </div>
    );
}
