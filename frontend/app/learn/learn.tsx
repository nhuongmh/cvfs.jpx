import React, { useState, useEffect } from "react";
import { MdCake, MdFlagCircle, MdHeartBroken, MdSnowboarding, MdOutlineRemoveRedEye  } from "react-icons/md";
import { Button } from "@nextui-org/react";
import { motion, AnimatePresence } from "framer-motion";
import { Card } from "models/card";
import { siteConfig } from "@/config/site";

const CardLearn: React.FC = () => {
    const [selectedLang] = useState("jpx");
    const [selectedGroup] = useState("minna");
    const [currentCard, setCurrentCard] = useState<Card | null>(null);
    const [error, setError] = useState<string>("");
    const [isShowBack, setIsShowBack] = useState<boolean>(false);

    const fetchCardFromServer = async () => {
        try {
            const url = `${siteConfig.server_url_prefix}/practice/${selectedLang}/${selectedGroup}/fetch`
            const response = await fetch(url);
            if (!response.ok) {
                const repjson = await response.json();
                throw new Error(`Server responded ${response.status}, ${repjson?.message}`);
            }
            const data = await response.json();
            setCurrentCard(data || null);
            console.log(currentCard)
            setError("");
        } catch (err: any) {
            setError("Failed to fetch card from server: " + err.message);
        }
    };

    const submitCardLearn = async (rating: string, fetchNext: boolean = true) => {
        if (!currentCard) return;
        try {
            const rateNum = convertRating(rating);
            const url = `${siteConfig.server_url_prefix}/practice/${selectedLang}/${selectedGroup}/submit?cardID=${currentCard.id}&rating=${rateNum}`;
            const response = await fetch(url, {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({}),
            });
            if (!response.ok) {
                const repjson = await response.json();
                throw new Error(`Server responded ${response.status}, ${repjson?.message}`);
            }
            setError("");
        } catch (err: any) {
            setError("Failed to submit card status: " + err.message);
        } finally {
            if (fetchNext) fetchCardFromServer();
        }
    };

    const convertRating = (rating: string) => {
        switch (rating) {
            case "again":
                return 1;
            case "hard":
                return 2;
            case "good":
                return 3;
            case "easy":
                return 4;
            default:
                throw new Error("Internal Error: Invalid rating");
        }
    };

    const toggleShowBack = () => {
        setIsShowBack(!isShowBack);
    };


    useEffect(() => {
        fetchCardFromServer();
    }, []);

    return (
        <div className="min-h-screen p-8">
            <div className="max-w-4xl mx-auto rounded-xl shadow-md overflow-hidden">
                <div className="p-8">
                    <AnimatePresence>
                        {currentCard && (
                            <motion.div
                                key={currentCard.id}
                                initial={{ opacity: 0, y: 20 }}
                                animate={{ opacity: 1, y: 0 }}
                                exit={{ opacity: 0, y: -20 }}
                                transition={{ duration: 0.3 }}
                                className="rounded-lg shadow-sm mb-6"
                            >
                                <div className="pb-11">
                                    <p className="text-4xl">{currentCard.front}</p>
                                </div>
                                {isShowBack && (
                                    <div>
                                        <h2 className="text-2xl font-bold mb-4 text-violet-700">Answer:</h2>
                                        <p className="text-4xl">{currentCard.back}</p>
                                    </div>
                                )}
                            </motion.div>
                        )}
                    </AnimatePresence>

                    {!isShowBack && (
                        <Button onPress={() => toggleShowBack()} className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded mx-2">
                            <MdOutlineRemoveRedEye  className="mr-2" /> Show Answer
                        </Button>
                    )}

                    {currentCard && isShowBack && (
                        <div className="flex justify-between items-center mb-6">
                            <div className="flex space-x-2">
                                <Button onPress={() => submitCardLearn("easy")} className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded mx-2">
                                    <MdCake className="mr-2" /> Easy
                                </Button>
                                <Button onPress={() => submitCardLearn("good")} className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded mx-2">
                                    <MdSnowboarding className="mr-2" /> Good
                                </Button>
                                <Button onPress={() => submitCardLearn("hard")} className="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded mx-2">
                                    <MdFlagCircle className="mr-2" /> Hard
                                </Button>
                                <Button onPress={() => submitCardLearn("again")} className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded mx-2">
                                    <MdHeartBroken className="mr-2" /> Again
                                </Button>
                            </div>
                        </div>
                    )}

                    {currentCard && (
                        <div className="mt-4">
                            <p className="text-lg">Due Review: <span className="text-xl">{new Date(currentCard.FsrsData.Due).toLocaleString()}</span></p>
                        </div>
                    )}

                    {error && (
                        <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-6" role="alert">
                            <p className="font-bold">Error</p>
                            <p>{error}</p>
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default CardLearn;
