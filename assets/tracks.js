import Map from 'ol/Map.js';
import View from 'ol/View.js';
import TileLayer from 'ol/layer/Tile.js';
import XYZ from 'ol/source/XYZ.js';
import VectorLayer from 'ol/layer/Vector.js';
import VectorSource from 'ol/source/Vector.js';
import Feature from 'ol/Feature.js';
import { LineString } from 'ol/geom.js';
import { fromLonLat } from 'ol/proj.js';
import Stroke from 'ol/style/Stroke.js';
import Style from 'ol/style/Style.js';

// Function to convert JSON points into OpenLayers LineString geometry
function createGPSLineString(data) {
  // Convert lat/lon data to map coordinates
  const coordinates = data.map(point => fromLonLat([point.Longitude, point.Latitude]));
  return new LineString(coordinates);
}

// get track id from current url. This assumes that we are running on the track
// view page. This assumption might not be true and this needs rework.
function getTrackId() {
  const url = window.location.href;
  const id = url.split('/').pop()
  return id;
}

// Function to fetch data from a JSON endpoint
async function fetchGPSTrack(map) {
  try {
    const response = await fetch('/track/' + getTrackId() + '/points'); // Replace with your JSON endpoint URL
    const jsonData = await response.json();

    // Create a LineString from the fetched data
    const trackLineString = createGPSLineString(jsonData);

    // Create a feature from the LineString
    const trackFeature = new Feature({
      geometry: trackLineString
    });

    // Create a vector source and layer to display the track
    const vectorSource = new VectorSource({
      features: [trackFeature]
    });

    const vectorLayer = new VectorLayer({
      source: vectorSource,
      style: new Style({
        stroke: new Stroke({
          color: 'blue',
          width: 3
        })
      })
    });

    // Add the vector layer to the map
    map.addLayer(vectorLayer);

    // Fit the map to the track
    map.getView().fit(vectorSource.getExtent(), { padding: [20, 20, 20, 20] });
  } catch (error) {
    console.error('Error fetching GPS track data:', error);
  }
}

document.addEventListener('DOMContentLoaded', function() {
  //confirm delete buttons
  //
  // Select all delete buttons by form class
  const deleteForms = document.querySelectorAll('.confirm');

  deleteForms.forEach(form => {
    form.addEventListener('submit', function (e) {
      if (!confirm('Are you sure you want to delete this object?')) {
        e.preventDefault(); // Prevent form submission if not confirmed
      }
    });
  });

  // close flash messages
  //
  //
  // Get all elements with data-close attribute
  const closeButtons = document.querySelectorAll('[data-close]');

  // Add click event to each close button
  closeButtons.forEach(button => {
    button.addEventListener('click', function () {
      // Get the id of the message to close
      const messageId = button.getAttribute('data-close');
      const messageElement = document.getElementById(messageId);

      // Hide the message element by adding the hidden class
      if (messageElement) {
        messageElement.classList.add('hidden');
      }
    });
  });

  const map = new Map({
    target: 'map',
    layers: [
      new TileLayer({
        title: 'OpenTopoMap',
        type: 'base',
        visible: true,
        source: new XYZ({
          url: 'https://{a-c}.tile.opentopomap.org/{z}/{x}/{y}.png'
        })
      }),
    ],
    view: new View({
      center: [0, 0],
      zoom: 2,
    }),
  });

  // Fetch and display the GPS track
  if (document.getElementById('map')) {
    fetchGPSTrack(map);
  }
}, false);


