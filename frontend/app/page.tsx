"use client";

import React from "react";
import { useRouter } from "next/navigation";
import VideoPlayer from "./components/VideoPlayer";
import apiClient from "../lib/api/client";

export default function FeedPage() {
  const router = useRouter();
  const [videos, setVideos] = React.useState<Array<{ id: string; video_url: string }> | null>(null);
  const [loading, setLoading] = React.useState(true);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    // Fetch For You feed from backend API
    apiClient.getForYouFeed()
      .then(response => {
        setVideos(response.videos || []);
        setLoading(false);
      })
      .catch(err => {
        console.error("Failed to fetch videos:", err);
        setError(err.message || "Failed to load videos");
        setLoading(false);

        // If unauthorized, redirect to login
        if (err.message?.includes("authorization") || err.message?.includes("token")) {
          setTimeout(() => router.push("/login"), 2000);
        }
      });
  }, [router]);

  if (loading) {
    return (
      <main className="h-screen w-full flex items-center justify-center">
        <div className="text-xl">Loading videos...</div>
      </main>
    );
  }

  if (error) {
    return (
      <main className="h-screen w-full flex items-center justify-center flex-col gap-4">
        <div className="text-xl text-red-500">Error: {error}</div>
        {(error.includes("authorization") || error.includes("token")) && (
          <div className="text-gray-600">Redirecting to login...</div>
        )}
      </main>
    );
  }

  if (!videos || videos.length === 0) {
    return (
      <main className="h-screen w-full flex items-center justify-center">
        <div className="text-xl">No videos available</div>
      </main>
    );
  }

  return (
    <main className="h-screen w-full snap-y snap-mandatory overflow-scroll">
      {videos.map(video => (
        <section key={video.id} className="snap-start h-screen flex items-center justify-center">
          <VideoPlayer src={video.video_url} />
        </section>
      ))}
    </main>
  );
}
