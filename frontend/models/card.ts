export interface Card {
    id: number;
    front: string;
    back: string;
    state: string;
    properties: {
        [key: string]: any;
    };
    FsrsData: FsrsData;
}

export interface FsrsData {
    Due: string;
    LastReview: string;
    ScheduledDays: number;
}