
export interface Article {
    id: number;
    title: string;
    content: string;
    origin: string;
    author: number;
    image: string;
    publish_date: string;
}

export interface ArticleReading {
    id : number;
    status : string;
    score : number;
    questions : Question[];
}

export interface Question {
    id: number;
    type: string;
    question: string;
    // user_answer_option: number;
    user_answer_str: string;
    options: string[];
    headings: string[];
    paragraph: string;
    correct_answer: string;
    correct: boolean;
}

export interface QuestionSubmit {
    id: number;
    user_answer_option: number;
    user_answer_str: string;
}

export interface ArticleTestSubmit {
    id: number;
    questions: QuestionSubmit[];
}
export interface QuestionResult {
    question_id: number;
    answer: string;
    user_answer: string;
    correct: boolean;
}

export interface ArticleTestResult {
    id: number;
    article_id: number;
    questions_result: QuestionResult[];
    score: number;
}

export interface ProposedWord {
    word: string;
    context_sentence: string;
    freq: number;
    ref_id: number;
    // selected: boolean;
}

export interface LearningWord {
    id: number;
    word: string;
    context_sentence: string;
    freq: number;
    ref_id: number;
}

export interface VocabList {
    id: number;
    name: string;
    vocabs: LearningWord[];
}