"use client";

import React, { useState } from 'react';
import { siteConfig } from "@/config/site";

export default function ContentForm() {
    const [formData, setFormData] = useState({
        origin: '',
        title: '',
        image: '',
        content: ''
    });
    const [isLoading, setIsLoading] = useState(false);

    const handleChange = (e: any) => {
        const { name, value } = e.target;
        setFormData(prevState => ({
            ...prevState,
            [name]: value
        }));
    };

    const handleLoadContent = async (e: any) => {
        e.preventDefault();
        if (!formData.origin) return;
        const url = `${siteConfig.private_url_prefix}/ie/article/url?url=${formData.origin}`;

        setIsLoading(true);
        try {
            const response = await fetch(url);
            if (!response.ok) {
                throw new Error('Failed to load content');
            }

            const data = await response.json();

            setFormData({
                ...formData,
                title: data.title || '',
                image: data.image || '',
                content: data.content || ''
            });
        } catch (error) {
            console.error('Error loading content:', error);
            // You could add error state handling here
        } finally {
            setIsLoading(false);
        }
    };

    const handleSubmit = async (e: any) => {
        e.preventDefault();

        const url = `${siteConfig.private_url_prefix}/ie/article`;
        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
            });

            if (!response.ok) {
                throw new Error('Failed to submit content');
            }

            // Handle successful submission
            alert('Content submitted successfully!');

            // Clear form
            setFormData({
                origin: '',
                title: '',
                image: '',
                content: ''
            });
        } catch (error) {
            console.error('Error submitting content:', error);
            // You could add error state handling here
        }
    };

    return (
        <div className="containerx px-8 mx-auto mt-16">
            <div className="max-w-4xl mx-auto p-6 bg-gray-30 rounded-lg">
                <h1 className="text-2xl font-semibold mb-6">New</h1>

                <form onSubmit={handleSubmit}>
                    <div className="mb-8">
                        <div className="flex gap-2">
                            <input
                                type="text"
                                name="origin"
                                placeholder="Paste link"
                                value={formData.origin}
                                onChange={handleChange}
                                className="bg-inherit flex-1 p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-200"
                            />
                            <button
                                type="button"
                                onClick={handleLoadContent}
                                disabled={isLoading}
                                className="bg-gray-900 text-white px-6 py-3 rounded-lg hover:bg-gray-800 transition duration-200"
                            >
                                {isLoading ? 'Loading...' : 'Load'}
                            </button>
                        </div>
                    </div>

                    <div className="border-t border-gray-200 my-6"></div>

                    <div className="mb-6">
                        <label className="block text-lg font-medium mb-2">
                            Title <span className="text-red-500">(*)</span>
                        </label>
                        <input
                            type="text"
                            name="title"
                            value={formData.title}
                            onChange={handleChange}
                            required
                            className="bg-inherit w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-200"
                        />
                    </div>

                    <div className="mb-6">
                        <label className="block text-lg font-medium mb-2">Cover Image</label>
                        <input
                            type="text"
                            name="image"
                            value={formData.image}
                            onChange={handleChange}
                            className="bg-inherit w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-200"
                            placeholder="Enter image URL"
                        />
                    </div>

                    <div className="mb-6">
                        <label className="block text-lg font-medium mb-2">Text</label>
                        <textarea
                            name="content"
                            value={formData.content}
                            onChange={handleChange}
                            className="bg-inherit w-full p-3 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-gray-200 min-h-[300px]"
                            placeholder="Paste content here...."
                        ></textarea>
                    </div>

                    <div className="flex justify-end">
                        <button
                            type="submit"
                            className="bg-gray-900 text-white px-6 py-3 rounded-lg hover:bg-gray-800 transition duration-200"
                        >
                            Submit
                        </button>
                    </div>
                </form>
            </div>
        </div>

    );
}