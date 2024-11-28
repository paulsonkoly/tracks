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

// Function to convert a segment of points into an OpenLayers LineString geometry
function createSegmentLineString(segment) {
  const coordinates = segment.map(point => fromLonLat([point.Longitude, point.Latitude]));
  return new LineString(coordinates);
}

// Function to fetch data from a JSON endpoint
async function fetchGPSTrack(map) {
  try {
    const response = await fetch(window.location.href + '/points', {headers: {Accept: "application/json"}}); 
    const jsonData = await response.json();

    // Iterate through each segment and create separate LineString features
    const features = jsonData.map(segment => {
      const lineString = createSegmentLineString(segment);
      return new Feature({
        geometry: lineString
      });
    });

    // Create a vector source and layer to display the segments
    const vectorSource = new VectorSource({
      features: features
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

    // Fit the map to the full extent of all segments
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


