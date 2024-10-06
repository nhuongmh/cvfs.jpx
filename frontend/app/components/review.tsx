import React, { useState, useEffect } from "react";
import { FaChevronLeft, FaChevronRight, FaSearch, FaFlag, FaCheck, FaTimes } from "react-icons/fa";
import { motion, AnimatePresence } from "framer-motion";

const SpacedRepetitionApp = () => {
  const [selectedDesk, setSelectedDesk] = useState("");
  const [currentCard, setCurrentCard] = useState(null);
  const [searchTerm, setSearchTerm] = useState("");
  const [error, setError] = useState("");
  const [note, setNote] = useState("");

  const dummyDesks = [
    { id: 1, name: "Mathematics" },
    { id: 2, name: "History" },
    { id: 3, name: "Science" },
    { id: 4, name: "Literature" },
    { id: 5, name: "Geography" },
  ];

  const dummyCards = [
    { id: 1, deskId: 1, content: "What is the Pythagorean theorem?" },
    { id: 2, deskId: 1, content: "Solve for x: 2x + 5 = 13" },
    { id: 3, deskId: 2, content: "When did World War II end?" },
    { id: 4, deskId: 2, content: "Who was the first President of the United States?" },
    { id: 5, deskId: 3, content: "What is photosynthesis?" },
  ];

  useEffect(() => {
    if (selectedDesk) {
      const cards = dummyCards.filter((card) => card.deskId === parseInt(selectedDesk));
      setCurrentCard(cards[0] || null);
      setError("");
    }
  }, [selectedDesk]);

  const filteredDesks = dummyDesks.filter((desk) =>
    desk.name.toLowerCase().includes(searchTerm.toLowerCase())
  );

  const handleDeskChange = (e) => {
    setSelectedDesk(e.target.value);
  };

  const handleSearchChange = (e) => {
    setSearchTerm(e.target.value);
  };

  const handleReviewStatus = (status) => {
    console.log(`Card marked as ${status}`);
    // Here you would typically update the card's status in your backend
  };

  const handleNoteChange = (e) => {
    setNote(e.target.value);
  };

  const handleSaveNote = () => {
    console.log("Note saved:", note);
    // Here you would typically save the note to your backend
    setNote("");
  };

  const handleNextCard = () => {
    const cards = dummyCards.filter((card) => card.deskId === parseInt(selectedDesk));
    const currentIndex = cards.findIndex((card) => card.id === currentCard.id);
    if (currentIndex < cards.length - 1) {
      setCurrentCard(cards[currentIndex + 1]);
    } else {
      setError("You've reached the end of this deck");
    }
  };

  const handlePreviousCard = () => {
    const cards = dummyCards.filter((card) => card.deskId === parseInt(selectedDesk));
    const currentIndex = cards.findIndex((card) => card.id === currentCard.id);
    if (currentIndex > 0) {
      setCurrentCard(cards[currentIndex - 1]);
    } else {
      setError("You're at the beginning of this deck");
    }
  };

  return (
    <div className="min-h-screen bg-gray-100 p-8">
      <div className="max-w-4xl mx-auto bg-white rounded-xl shadow-md overflow-hidden">
        <div className="p-8">
          <h1 className="text-3xl font-bold mb-6 text-center text-indigo-600">Spaced Repetition Review App</h1>
          
          {/* Desk Selection */}
          <div className="mb-6">
            <label htmlFor="desk-select" className="block text-sm font-medium text-gray-700 mb-2">
              Select a Desk
            </label>
            <div className="relative">
              <FaSearch className="absolute left-3 top-3 text-gray-400" />
              <input
                type="text"
                placeholder="Search desks..."
                value={searchTerm}
                onChange={handleSearchChange}
                className="pl-10 pr-4 py-2 w-full border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
              />
            </div>
            <select
              id="desk-select"
              value={selectedDesk}
              onChange={handleDeskChange}
              className="mt-2 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
              aria-label="Select a desk"
            >
              <option value="">Choose a desk</option>
              {filteredDesks.map((desk) => (
                <option key={desk.id} value={desk.id}>
                  {desk.name}
                </option>
              ))}
            </select>
          </div>

          {/* Card Display */}
          <AnimatePresence>
            {currentCard && (
              <motion.div
                key={currentCard.id}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                exit={{ opacity: 0, y: -20 }}
                transition={{ duration: 0.3 }}
                className="bg-gray-50 p-6 rounded-lg shadow-sm mb-6"
              >
                <h2 className="text-xl font-semibold mb-4">Card Content</h2>
                <p className="text-gray-700">{currentCard.content}</p>
              </motion.div>
            )}
          </AnimatePresence>

          {/* Navigation and Review Status */}
          {currentCard && (
            <div className="flex justify-between items-center mb-6">
              <button
                onClick={handlePreviousCard}
                className="flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                aria-label="Previous card"
              >
                <FaChevronLeft className="mr-2" /> Previous
              </button>
              <div className="flex space-x-2">
                <button
                  onClick={() => handleReviewStatus("easy")}
                  className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-green-700 bg-green-100 hover:bg-green-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500"
                  aria-label="Mark as easy"
                >
                  Easy
                </button>
                <button
                  onClick={() => handleReviewStatus("moderate")}
                  className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-yellow-700 bg-yellow-100 hover:bg-yellow-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-yellow-500"
                  aria-label="Mark as moderate"
                >
                  Moderate
                </button>
                <button
                  onClick={() => handleReviewStatus("difficult")}
                  className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-red-700 bg-red-100 hover:bg-red-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
                  aria-label="Mark as difficult"
                >
                  Difficult
                </button>
              </div>
              <button
                onClick={handleNextCard}
                className="flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                aria-label="Next card"
              >
                Next <FaChevronRight className="ml-2" />
              </button>
            </div>
          )}

          {/* Note Input */}
          {currentCard && (
            <div className="mb-6">
              <label htmlFor="note-input" className="block text-sm font-medium text-gray-700 mb-2">
                Add a Note
              </label>
              <div className="flex">
                <input
                  id="note-input"
                  type="text"
                  value={note}
                  onChange={handleNoteChange}
                  className="flex-grow mr-2 px-4 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
                  placeholder="Type your note here..."
                />
                <button
                  onClick={handleSaveNote}
                  className="px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  Save Note
                </button>
              </div>
            </div>
          )}

          {/* Error Handling */}
          {error && (
            <div className="bg-red-100 border-l-4 border-red-500 text-red-700 p-4 mb-6" role="alert">
              <p className="font-bold">Error</p>
              <p>{error}</p>
            </div>
          )}

          {/* Additional Actions */}
          {currentCard && (
            <div className="flex justify-end space-x-2">
              <button
                onClick={() => console.log("Card flagged for review")}
                className="flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-gray-700 bg-gray-100 hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-gray-500"
                aria-label="Flag for review"
              >
                <FaFlag className="mr-2" /> Flag for Review
              </button>
              <button
                onClick={() => console.log("Card marked as reviewed")}
                className="flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                aria-label="Mark as reviewed"
              >
                <FaCheck className="mr-2" /> Mark as Reviewed
              </button>
            </div>
          )}
        </div>
      </div>
    </div>
  );
};

export default SpacedRepetitionApp;
