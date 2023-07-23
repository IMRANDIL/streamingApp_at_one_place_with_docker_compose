import React, { useState, useEffect } from "react";
import axios from "axios";
import { toast } from "react-toastify";
import "./App.css";
import VideoList from "./component/VideoList";
export default function App() {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [videoID, setVideoID] = useState<string | null>(null);
  const [videos, setVideos] = useState<Video[]>([]);
  const [loading, setLoading] = useState(false); // State for the loading spinner
  const [currentPlayingVideo, setCurrentPlayingVideo] = useState<string | null>(
    null
  );

  const isVideoFileType = (filename: string) => {
    // Add more video formats here if needed
    const videoFormats = [".mp4", ".avi"];
    const ext = filename.slice(filename.lastIndexOf(".")).toLowerCase();
    return videoFormats.includes(ext);
  };

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files.length > 0) {
      const file = e.target.files[0];
      if (file.size <= 200 * 1024 * 1024) {
        // Check if the file size is within the allowed limit (10MB)
        if (isVideoFileType(file.name)) {
          setSelectedFile(file);
        } else {
          toast.warn("Invalid file type. Only mp4 and avi files are allowed.");
        }
      } else {
        toast.warn("File size exceeds the maximum limit of 10MB.");
      }
    }
  };

  const handleUpload = async () => {
    if (!selectedFile) return;

    const formData = new FormData();
    formData.append("video", selectedFile);

    try {
      setLoading(true); // Show loading spinner while uploading
      const response: any = await axios.post(
        "http://localhost:8080/upload",
        formData,
        {
          headers: {
            "Content-Type": "multipart/form-data",
          },
        }
      );

      // Video uploaded successfully, get the video ID from the response
      setVideoID(response?.data.id);

      // Fetch the updated list of videos
      fetchVideos();
    } catch (error: any) {
      toast.error(error?.response?.data);
    } finally {
      setLoading(false); // Hide loading spinner after upload completes
    }
  };

  const fetchVideos = async () => {
    try {
      const response: any = await axios.get("http://localhost:8080/videos");
      setVideos(response.data);
    } catch (error: any) {
      toast.error(error?.response?.data);
    }
  };

  useEffect(() => {
    // Fetch the initial list of videos on component mount
    fetchVideos();
  }, []);

  const handleVideoPlay = (videoID: string) => {
    // If the clicked video is the same as the currently playing video, pause it
    if (currentPlayingVideo === videoID) {
      setCurrentPlayingVideo(null);
    } else {
      // If a different video is clicked, set the video ID and pause the currently playing video (if any)
      setVideoID(videoID);
      setCurrentPlayingVideo(videoID);
    }
  };

  const handleVideoPause = () => {
    setCurrentPlayingVideo(null);
  };

  return (
    <div className="grandParent">
      <div className="upload">
        {loading ? (
          // Show the loading spinner if the file is being uploaded
          <div className="spinner" />
        ) : (
          <img
            src="https://media.giphy.com/media/7T8uLiWEsnytWkQKvF/giphy.gif"
            alt="Video Upload GIF"
            draggable="false"
          />
        )}
        <label htmlFor="fileInput">Choose a video</label>
        <input type="file" id="fileInput" onChange={handleFileChange} />
        <button onClick={handleUpload}>Upload</button>
      </div>
      <h2 style={{ textAlign: "center" }}>All Videos</h2>

      <VideoList
        videos={videos}
        videoID={videoID}
        currentPlayingVideo={currentPlayingVideo}
        handleVideoPlay={handleVideoPlay}
        handleVideoPause={handleVideoPause}
      />
    </div>
  );
}

interface Video {
  id: string;
  name: string;
  url: string;
}
