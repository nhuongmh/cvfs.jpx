"use client";
import React, { useEffect, useState } from 'react';
import { siteConfig } from "@/config/site";
import { Table, TableHeader, TableBody, TableColumn, TableRow, TableCell } from "@heroui/table";
import Link from "next/link";

interface ArticleData {
    id: number;
    title: string;
    created_at: string;
    status: string;
    origin: string;
}

const renderCell = (article: ArticleData, columnKey: React.Key) => {
    const cellValue = article[columnKey as keyof ArticleData];

    switch (columnKey) {
        case "title":
            return (
                <Link href={`/reading/${article.id}`} rel="noopener noreferrer">
                    {cellValue}
                </Link>
            );
        case "created_at":
            return new Date(cellValue).toLocaleString(); // Format the date
        default:
            return cellValue;
    }
};

export default function ArticlesPage() {
    const [articles, setArticles] = useState<ArticleData[]>([]);
    const [loading, setLoading] = useState<boolean>(true);

    const columns = [
        {
            key: "id",
            label: "ID",
        },
        {
            key: "title",
            label: "TITLE",
        },
        {
            key: "created_at",
            label: "CREATED AT",
        },
        {
            key: "status",
            label: "STATUS",
        },
    ];

    useEffect(() => {
        const fetchArticles = async () => {
            try {
                const response = await fetch(`${siteConfig.private_url_prefix}/ie/article`);
                const data = await response.json();
                setArticles(data.articles);
            } catch (error) {
                console.error('Error fetching articles:', error);
            } finally {
                setLoading(false);
            }
        };

        fetchArticles();
    }, []);

    if (loading) {
        return <div>Loading...</div>;
    }

    return (
        <div className="container px-8 mx-auto mt-16 lg:mt-32">
            <Table aria-label="Articles Table">
                <TableHeader columns={columns}>
                    {(column) => <TableColumn key={column.key}>{column.label}</TableColumn>}
                </TableHeader>
                <TableBody items={articles}>
                    {(item) => (
                        <TableRow key={item.id}>
                            {(columnKey) => <TableCell>{renderCell(item, columnKey)}</TableCell>}
                        </TableRow>
                    )}
                </TableBody>
            </Table>
        </div>
    );
};
