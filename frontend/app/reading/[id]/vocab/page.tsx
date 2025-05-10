"use client";
import React, { useState, useEffect } from "react";
import { Article, LearningWord } from '@/models/article';
import { Button, Checkbox, Input, CheckboxGroup } from "@heroui/react";
import { Table, TableHeader, TableColumn, TableBody, TableRow, TableCell, getKeyValue } from "@heroui/react";
import type {Selection} from "@heroui/react";
import { ProposedWord } from "@/models/article";
import { useSearchParams, usePathname } from "next/navigation";
import { siteConfig } from "@/config/site";
import '../ielts.css';

const ProposalVocabPage = ({
    articleId,
}: {
    articleId: string;
}) => {
    const [proposalVocabs, setProposalVocabs] = useState<ProposedWord[]>([]);
    const [selectedVocabs, setSelectedVocabs] = useState<Selection>(new Set([]));
    const [article, setArticle] = useState<Article | null>(null);
    const [newWord, setNewWord] = useState("");
    const [error, setError] = useState<string | null>(null);

    const columns = [
        {
          key: "word",
          label: "VOCAB",
        },
        {
          key: "context",
          label: "CONTEXT",
        },
        {
          key: "freq",
          label: "FREQ",
        },
      ];

    useEffect(() => {
        const fetchVocabs = async () => {
            try {
                const url = `${siteConfig.private_url_prefix}/ie/article/${articleId}/proposed_vocab`;
                const response = await fetch(url);
                if (!response.ok) {
                    throw new Error("Failed to fetch article");
                }
                const data = await response.json();
                setProposalVocabs(data);
            } catch (error) {
                console.error("Error fetching vocabularies:", error);
            }
        };
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
        fetchVocabs();
        console.log("fetchVocabs", proposalVocabs);
    }, [articleId]);

    // const handleCheckboxChange = (words: string[]) => {
    //     setProposalVocabs((prev) =>
    //         prev.map((vocab) => ({
    //         ...vocab,
    //         selected: words.includes(vocab.word),
    //         }))
    //     );
    // };

    const handleAddWord = () => {
        if (newWord.trim()) {
            setProposalVocabs((prev) => [
                ...prev,
                { word: newWord, context_sentence: "", freq: 0 } as ProposedWord,
            ]);
            setNewWord("");
        }
    };

    const formatContent = (content: string) => {
        if (!content) return [];
        return content.split('\n').map((paragraph, index) =>
            <p key={index}>{paragraph}</p>
        );
    };

    const handleSubmit = async () => {
        const selected = proposalVocabs.filter((vocab) =>
            selectedVocabs === "all" || selectedVocabs.has(vocab.word)
        );
        try {
            await fetch(
                `${siteConfig.private_url_prefix}/ie/article/${articleId}/proposed_vocab`,
                {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify(Object.values(selected)),
                }
            );
            alert("Vocabularies submitted successfully!");
        } catch (error) {
            console.error("Error submitting vocabularies:", error);
        }
    };

    if (error) {
        return (
            <div className="error-container">
                <h2>Error</h2>
                <p>{error}</p>
                <button onClick={() => window.location.reload()}>Try Again</button>
            </div>
        );
    }

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
                        <Table
                            aria-label="Controlled table example with dynamic content"
                            selectionMode="multiple"
                            selectedKeys={selectedVocabs}
                            onSelectionChange={setSelectedVocabs}
                        >
                            <TableHeader columns={columns}>
                                {(column) => <TableColumn key={column.key}>{column.label}</TableColumn>}
                            </TableHeader>
                            <TableBody items={proposalVocabs}>
                                {(item) => (
                                    <TableRow key={item.word}>
                                        <TableCell>{item.word}</TableCell>
                                        <TableCell>{item.context_sentence}</TableCell>
                                        <TableCell>{item.freq}</TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                        <div className="mt-4">
                            <Input
                                value={newWord}
                                onChange={(e) => setNewWord(e.target.value)}
                                placeholder="Add new vocabulary"
                            />
                            <Button onPress={handleAddWord} className="ml-2">
                                Add
                            </Button>
                        </div>
                        <Button onPress={handleSubmit} className="mt-4">
                            Submit
                        </Button>
                    </div>

                </div>
            </div>
        </div>

    );
};

const LearningVocabPage = ({
    learningVocabs,
}: {
    learningVocabs: LearningWord[];
}) => {
    return (
        <div className="p-4">
            <h1 className="text-xl font-bold mb-4">Learning Vocabularies</h1>
            <CheckboxGroup>
                {learningVocabs?.map((vocab, index) => (
                    <div key={index} className="flex items-center mb-2">
                        <div className="ml-2">
                            <p className="font-medium">{vocab.word}</p>
                            <p className="text-sm text-gray-500">{vocab.context}</p>
                        </div>
                    </div>
                ))}
            </CheckboxGroup>
        </div>
    );
};

const VocabPage = () => {
    const [learningVocabs, setLearningVocabs] = useState<LearningWord[] | null>(
        null
    );
    const [isLoading, setIsLoading] = useState(true);

    const searchParams = useSearchParams();
    const pathname = usePathname();
    if (pathname === null || searchParams === null) {
        return null;
    }
    const articleId = pathname.split("/")[2] || "";

    useEffect(() => {
        const fetchVocabs = async () => {
            try {
                const url = `${siteConfig.private_url_prefix}/ie/article/${articleId}/vocab`;
                const response = await fetch(url);
                if (response.ok) {
                    const data = await response.json();
                    setLearningVocabs(data);
                } else {
                    setLearningVocabs(null);
                }
            } catch (error) {
                console.error("Error fetching vocabularies:", error);
                setLearningVocabs(null);
            } finally {
                setIsLoading(false);
            }
        };
        fetchVocabs();
    }, [articleId]);

    if (isLoading) {
        if (isLoading) {
            return (
                <div className="loading-container">
                    <div className="loading-spinner"></div>
                    <p>Loading test content...</p>
                </div>
            );
        }
    }

    return learningVocabs ? (
        <LearningVocabPage learningVocabs={learningVocabs} />
    ) : (
        <ProposalVocabPage articleId={articleId} />
    );
};

export default VocabPage;