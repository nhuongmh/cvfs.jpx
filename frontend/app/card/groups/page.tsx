"use client";
import React, { useEffect, useState } from 'react';
import { siteConfig } from "@/config/site";
import { Table, TableHeader, TableBody, TableColumn, TableRow, TableCell } from "@heroui/table";
import Link from "next/link";


interface GroupData {
    group: string;
    num_cards: number;
    proposal: number;
    learning: number;
    discard: number;
    save: number;
}

const renderCell = (group: GroupData, columnKey: React.Key) => {
    const cellValue = group[columnKey as keyof GroupData];

    switch (columnKey) {
        case "group":
            return (
                <Link
                    href={`/card/groups/${group.group}`}
                >
                    {cellValue}
                </Link>
            );
        case "proposal":
            return (
                <Link
                    href={`/card/groups/process/${group.group}`}
                >
                    {cellValue}
                </Link>
            );
        case "learning":
            return (
                <Link
                    href={`/card/groups/learn/${group.group}`}
                >
                    {cellValue}
                </Link>
            );
        default:
            return cellValue;
    }
};

export default function GroupsPage() {
    const [groups, setGroups] = useState<GroupData[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    const columns = [
        {
            key: "group",
            label: "GROUP",
        },
        {
            key: "num_cards",
            label: "NUMBER OF CARDS",
        },
        {
            key: "proposal",
            label: "PROPOSAL",
        },
        {
            key: "learning",
            label: "LEARNING",
        },
        {
            key: "discard",
            label: "DISCARD",
        },
        {
            key: "save",
            label: "SAVE",
        },
    ];

    useEffect(() => {
        const fetchGroups = async () => {
            try {
                const response = await fetch(`${siteConfig.server_url_prefix}/practice/jp/stats`);
                const data: GroupData[] = await response.json();
                setGroups(data);
            } catch (error) {
                console.error('Error fetching group data:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchGroups();
    }, []);

    if (loading) {
        return <div>Loading...</div>;
    }

    return (
        <div className="container px-8 mx-auto mt-16 lg:mt-32 ">
            <Table aria-label="Example table with dynamic content">
                <TableHeader columns={columns}>
                    {(column) => <TableColumn key={column.key}>{column.label}</TableColumn>}
                </TableHeader>
                <TableBody items={groups}>
                    {(item) => (
                        <TableRow key={item.group}>
                            {(columnKey) => <TableCell>{renderCell(item, columnKey)}</TableCell>}
                        </TableRow>
                    )}
                </TableBody>
            </Table>
        </div>
    );
};
