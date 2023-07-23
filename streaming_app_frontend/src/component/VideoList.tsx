import React from "react";
import ReactPlayer from "react-player";

interface Video {
  id: string;
  name: string;
  url: string;
}

interface VideoListProps {
  videos: Video[];
  videoID: string | null;
  currentPlayingVideo: string | null;
  handleVideoPlay: (videoID: string) => void;
  handleVideoPause: () => void;
}

const VideoList: React.FC<VideoListProps> = ({
  videos,
  currentPlayingVideo,
  handleVideoPlay,
  handleVideoPause,
}) => {
  return (
    <div className="video-list">
      {videos.length > 0 ? (
        <div className="video-list-inner">
          {videos.map((video) => (
            <div key={video.id} className="video-card">
              <p>Video ID: {video.id}</p>
              <p>Video Name: {video.name}</p>
              <ReactPlayer
                url={video.url}
                width="320"
                height="240"
                controls
                playing={currentPlayingVideo === video.id}
                onPlay={() => handleVideoPlay(video.id)}
                onPause={handleVideoPause}
              />
            </div>
          ))}
        </div>
      ) : (
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            marginTop: "20px",
          }}
        >
          <img
            src="https://media.giphy.com/media/fipN1GOuDK8txSqay3/giphy.gif"
            alt="No Videos GIF"
            draggable="false"
          />
        </div>
      )}
    </div>
  );
};

export default VideoList;
