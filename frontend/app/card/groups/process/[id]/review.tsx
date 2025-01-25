import React, { useState, useEffect } from "react";
import { FaCheck, FaTimes, FaArchive, FaEdit } from "react-icons/fa";
import { Modal, ModalContent, ModalHeader, ModalBody, ModalFooter, Button, useDisclosure, Input } from "@nextui-org/react";
import { motion, AnimatePresence } from "framer-motion";
import { Card } from "models/card";
import { siteConfig } from "@/config/site";

interface ProposalCheckProps {
  lang?: string;
  group: string;
}

const ProposalCheck: React.FC<ProposalCheckProps> = ({lang="jpx", group}) => {
  const [currentCard, setCurrentCard] = useState<Card | null>(null);
  const [currentCardEdit, setCurrentCardEdit] = useState<Card | null>(null);
  const [error, setError] = useState<string>("");
  const [note, setNote] = useState<string>("");
  const { isOpen, onOpen, onOpenChange } = useDisclosure();

  const fetchCardFromServer = async () => {
    try {
      if (!group) {
        setError("No group selected");
        return;
      };
      const response = await fetch(`${siteConfig.server_url_prefix}/process/${lang}/fetch?group=${group}`);
      if (!response.ok) {
        const repjson = await response.json();
        throw new Error(`Server responded ${response.status}, ${repjson?.message}`);
      }
      const data = await response.json();
      updateCurrentCardSync(data || null);
      console.log(currentCard)
      setError("");
    } catch (err: any) {
      setError("Failed to fetch card from server: " + err.message);
    }
  };

  const submitCardStatus = async (status: string, fetchNext: boolean = true) => {
    if (!currentCard) return;
    try {
      const url = `${siteConfig.server_url_prefix}/process/${lang}/submit?cardID=${currentCard.id}&status=${status}`;
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

  const updateCardEdit = async () => {
    if (!currentCardEdit) return;
    try {
      console.log("updating card: ")
      const url = `${siteConfig.server_url_prefix}/process/jpx/edit`;
      const response = await fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(currentCardEdit),
      });
      if (!response.ok) throw new Error("Network response was not ok");
      const data = await response.json();
      updateCurrentCardSync(data || null);
      setError("");
    } catch (err: any) {
      setError("Failed to edit card: " + err.message);
    }
  };

  const updateCurrentCardSync = async (card: Card | null) => {
    setCurrentCard(card);
    setCurrentCardEdit(card);
  }

  const syncCardEdit = async () => {
    setCurrentCardEdit(currentCard);
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
                className="p-6 rounded-lg shadow-sm mb-6"
              >
                <div className="pb-11">
                  <h2 className="text-2xl font-bold mb-4 text-fuchsia-700">Front:</h2>
                  <p className="text-3xl">{currentCard.front}</p>
                </div>
                <div>
                  <h2 className="text-2xl font-bold mb-4 text-violet-700">Back:</h2>
                  <p className="text-3xl">{currentCard.back}</p>
                </div>
                {/* display cards properties */}
                <div>
                  <h2 className="text-2xl font-bold mb-4 text-indigo-700">Properties:</h2>
                  <ul>
                    {Object.entries(currentCard.properties).map(([key, value]) => (
                      <li key={key} className="mb-2">
                        <strong>{key}:</strong> {value}
                      </li>
                    ))}
                  </ul>
                </div>
              </motion.div>
            )}
          </AnimatePresence>

          {currentCard && (
            <div className="flex justify-between items-center mb-6">
              <div className="flex space-x-2">
                <Button onPress={() => submitCardStatus("learn")} className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded mx-2">
                  <FaCheck className="mr-2" /> Learn
                </Button>
                <Button onPress={() => submitCardStatus("discard")} className="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded mx-2">
                  <FaTimes className="mr-2" /> Discard
                </Button>
                <Button onPress={() => submitCardStatus("save")} className="bg-yellow-500 hover:bg-yellow-700 text-white font-bold py-2 px-4 rounded mx-2">
                  <FaArchive className="mr-2" /> Archive
                </Button>
                <Button onPress={onOpen} className="bg-gray-500 hover:bg-gray-700 text-white font-bold py-2 px-4 rounded mx-2">
                  <FaEdit className="mr-2" /> Edit
                </Button>
              </div>
              <Modal isOpen={isOpen} onOpenChange={onOpenChange} placement="top-center">
                <ModalContent>
                  {(onClose) => (
                    <>
                      <ModalHeader className="flex flex-col gap-1">Edit</ModalHeader>
                      <ModalBody>
                        <Input
                          autoFocus
                          placeholder="Front"
                          variant="bordered"
                          value={currentCardEdit?.front || ""}
                          onChange={(e) => setCurrentCardEdit({ ...currentCardEdit, front: e.target.value } as Card)}
                        />
                        <Input
                          placeholder="Back"
                          variant="bordered"
                          value={currentCardEdit?.back || ""}
                          onChange={(e) => setCurrentCardEdit({ ...currentCardEdit, back: e.target.value } as Card)}
                        />
                      </ModalBody>
                      <ModalFooter>
                        <Button color="danger" variant="flat" onPress={() => { onClose(); syncCardEdit(); }}>Close</Button>
                        <Button color="primary" onPress={() => { onClose(); updateCardEdit(); }}>Save</Button>
                      </ModalFooter>
                    </>
                  )}
                </ModalContent>
              </Modal>
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

export default ProposalCheck;
