import Map from 'ol/Map.js';
import View from 'ol/View.js';
import TileLayer from 'ol/layer/Tile.js';
import XYZ from 'ol/source/XYZ.js';
import VectorLayer from 'ol/layer/Vector.js';
import VectorSource from 'ol/source/Vector.js';
import Feature from 'ol/Feature.js';
import { extend as extendExtent, createEmpty as createEmptyExtent } from 'ol/extent';
import { LineString } from 'ol/geom.js';
import { fromLonLat } from 'ol/proj.js';
import Stroke from 'ol/style/Stroke.js';
import Style from 'ol/style/Style.js';

// Define a palette of high-contrast colors
const highContrastColors = [
  "#FF5733", // Bright orange-red
  "#33FF57", // Bright green
  "#3357FF", // Bright blue
  "#FF33A1", // Vibrant pink
  "#FFD700", // Golden yellow
  "#00FFFF", // Cyan
  "#FF4500", // Orange-red
  "#8A2BE2", // Blue-violet
  "#00FF7F", // Spring green
  "#DC143C", // Crimson
];

// Function to get a color from the palette
function getHighContrastColor(index) {
  return highContrastColors[index % highContrastColors.length];
}
// Function to create a table row for a track
function createTrackTableRow(track, color) {
  const row = document.createElement("tr");
  row.className = "border-b border-gray-200";

  const cell = document.createElement("td");
  cell.className = "py-2 px-4 flex items-center justify-between"; // Align content horizontally

  // Track name with a link
  const trackLink = document.createElement("a");
  trackLink.href = `/track/${track.id}`;
  trackLink.textContent = track.name;
  trackLink.className = "text-green-600 hover:underline";

  // Inline colored line (styled `hr`)
  const colorLine = document.createElement("hr");
  colorLine.style.border = `2px solid ${color}`; // Set the color
  colorLine.style.width = "50px"; // Set a shorter width
  colorLine.style.marginLeft = "10px"; // Add some space between the name and the line

  // Add both elements to the cell
  cell.appendChild(trackLink);
  cell.appendChild(colorLine);

  // Add the cell to the row
  row.appendChild(cell);

  return row;
}

// Function to convert a segment of points into an OpenLayers LineString geometry
function createSegmentLineString(segment) {
  const coordinates = segment.map(point => fromLonLat([point.Longitude, point.Latitude]));
  return new LineString(coordinates);
}

// Function to fetch and render a single track and update the overall extent
async function fetchAndRenderTrack(map, track, color, combinedExtent) {
  try {
    const response = await fetch(`/track/${track.id}/points`, { headers: { Accept: "application/json" } });
    const jsonData = await response.json();

    // Create LineString features for the track segments
    const features = jsonData.map(segment => {
      const lineString = createSegmentLineString(segment);
      return new Feature({
        geometry: lineString,
      });
    });

    // Create a vector source and layer for the track
    const vectorSource = new VectorSource({ features: features });

    const vectorLayer = new VectorLayer({
      source: vectorSource,
      style: new Style({
        stroke: new Stroke({
          color: color,
          width: 3,
        }),
      }),
    });

    // Add the vector layer to the map
    map.addLayer(vectorLayer);

    // Extend the combined extent to include this track's extent
    const trackExtent = vectorSource.getExtent();
    extendExtent(combinedExtent, trackExtent);
  } catch (error) {
    console.error(`Error fetching data for track ID ${track.id}:`, error);
  }
}

// Function to fetch and display the GPS tracks and populate the table
async function fetchGPSTrack(map) {
  try {
    const currentUrl = window.location.href;
    const combinedExtent = createEmptyExtent(); // Create an empty extent to combine all tracks

    if (currentUrl.includes("/collection/")) {
      // Fetch the collection data
      const response = await fetch(currentUrl, { headers: { Accept: "application/json" } });
      const collectionData = await response.json();

      // Extract tracks from the collection object
      const tracks = collectionData.tracks;

      // Get the table body to populate rows
      const tableBody = document.querySelector("table tbody");
      if (!tableBody) {
        console.error("Track table body not found!");
        return;
      }

      // Fetch and render each track
      tracks.forEach(async (track, index) => {
        const color = getHighContrastColor(index);

        // Add a row to the table for this track
        const row = createTrackTableRow(track, color);
        tableBody.appendChild(row);

        // Render the track on the map and update the combined extent
        await fetchAndRenderTrack(map, track, color, combinedExtent);

        // Fit the map to the combined extent of all tracks
        map.getView().fit(combinedExtent, { padding: [20, 20, 20, 20] });
      });
    } else if (currentUrl.includes("/track/")) {
      // For individual tracks, fetch and render the track directly
      const trackId = currentUrl.split("/track/")[1].split("/")[0];
      const color = highContrastColors[0]; // Default color for individual tracks
      await fetchAndRenderTrack(map, { id: trackId }, color, combinedExtent);

      // Fit the map to the extent of the single track
      map.getView().fit(combinedExtent, { padding: [20, 20, 20, 20] });
    }
  } catch (error) {
    console.error("Error fetching GPS track data:", error);
  }
}

document.addEventListener("DOMContentLoaded", function () {
  // Confirm delete buttons
  const deleteForms = document.querySelectorAll(".confirm");
  deleteForms.forEach(form => {
    form.addEventListener("submit", function (e) {
      if (!confirm("Are you sure you want to delete this object?")) {
        e.preventDefault(); // Prevent form submission if not confirmed
      }
    });
  });

  // Close flash messages
  const closeButtons = document.querySelectorAll("[data-close]");
  closeButtons.forEach(button => {
    button.addEventListener("click", function () {
      const messageId = button.getAttribute("data-close");
      const messageElement = document.getElementById(messageId);
      if (messageElement) {
        messageElement.classList.add("hidden");
      }
    });
  });

  // Initialize the map
  const map = new Map({
    target: "map",
    layers: [
      new TileLayer({
        title: "OpenTopoMap",
        type: "base",
        visible: true,
        source: new XYZ({
          url: "https://{a-c}.tile.opentopomap.org/{z}/{x}/{y}.png",
        }),
      }),
    ],
    view: new View({
      center: [0, 0],
      zoom: 2,
    }),
  });

  // Fetch and display the GPS tracks or collection
  if (document.getElementById("map")) {
    fetchGPSTrack(map);
  }
}, false);

