"use client";
import React, { useState, useEffect } from "react";
import { Article, LearningWord, VocabList } from '@/models/article';
import { Table, TableHeader, TableColumn, TableBody, TableRow, TableCell, getKeyValue } from "@heroui/react";
import { siteConfig } from "@/config/site";

const LearningVocabPage = ({
    vocabList,
    articleId

}: {
    vocabList: VocabList;
    articleId: string;
}) => {
    const [article, setArticle] = useState<Article | null>(null);
    const [error, setError] = useState<string | null>(null);

    useEffect(() => {
        const fetchArticle = async () => {
            try {
                const url = `${siteConfig.private_url_prefix}/ie/article/${articleId}`;
                const response = await fetch(url);
                if (!response.ok) {
                    throw new Error('Failed to fetch article');
                }
                const data = await response.json();
                setArticle(data);
            } catch (err) {
                setError('Error fetching article: ' + (err as Error).message);
                console.error('Error fetching article:', err);
            }
        };
        fetchArticle();
    }, [articleId]);

    const formatContent = (content: string) => {
        if (!content) return [];
        return content.split('\n').map((paragraph, index) =>
            <p key={index}>{paragraph}</p>
        );
    };

    return (
        <div className="app-container">
            <div className="main-content">
                <div className="reading-section">

                    <div className="passage-container">
                        {article && (
                            <>
                                <div className="passage-header">
                                    <div className="passage-logo">
                                    </div>
                                    <div className="passage-title">
                                        <h2>{article.title}</h2>
                                    </div>
                                </div>

                                <div className="passage-content">
                                    <h2>{article.title}</h2>
                                    {article.image && (
                                        <div className="passage-image">
                                            <img src={article.image} alt={article.title} />
                                        </div>
                                    )}
                                    <div className="passage-text">
                                        {formatContent(article.content)}
                                    </div>
                                </div>
                            </>
                        )}
                    </div>

                </div>

                <div className="questions-section">
                    <div className="p-4">
                        <h1 className="text-xl font-bold mb-4">Proposed Vocabularies</h1>
                        {Array.isArray(vocabList.vocabs) ? (
                            <Table>
                                <TableHeader>
                                    <TableColumn>Word</TableColumn>
                                    <TableColumn>Context</TableColumn>
                                </TableHeader>
                                <TableBody items={vocabList.vocabs}>
                                    {(item) => (
                                        <TableRow key={item.word}>
                                            <TableCell>{item.word}</TableCell>
                                            <TableCell>{item.context_sentence}</TableCell>
                                        </TableRow>
                                    )}
                                </TableBody>
                            </Table>
                        ) : (
                            <div className="error-message">
                                <p>Error: Learning vocabularies data is not available or not iterable.</p>
                            </div>
                        )}
                    </div>

                </div>
            </div>
        </div>
    );
};

export default LearningVocabPage;