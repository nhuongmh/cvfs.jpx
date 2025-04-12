
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
    user_answer_option: number;
    user_answer_str: string;
    options: string[];
    headings: string[];
    paragraph: string;
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