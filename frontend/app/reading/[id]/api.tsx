// src/services/api.js

const API_BASE_URL = 'http://localhost:9000/private/api/v1/ie';

export const fetchArticle = async (articleId: number) => {
  try {
    const response = await fetch(`${API_BASE_URL}/article/${articleId}/`);
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching article:', error);
    throw error;
  }
};

export const fetchQuestions = async (articleId: number) => {
  try {
    const response = await fetch(`${API_BASE_URL}/article/${articleId}/questions`);
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error fetching questions:', error);
    throw error;
  }
};

export const submitAnswers = async (articleId: number, answers: string) => {
  try {
    const response = await fetch(`${API_BASE_URL}/article/${articleId}/answers`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ answers }),
    });
    if (!response.ok) {
      throw new Error(`API error: ${response.status}`);
    }
    return await response.json();
  } catch (error) {
    console.error('Error submitting answers:', error);
    throw error;
  }
};