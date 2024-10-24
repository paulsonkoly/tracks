import Map from 'ol/Map.js';
import OSM from 'ol/source/OSM.js';
import XYZ from 'ol/source/XYZ.js';
import TileLayer from 'ol/layer/Tile.js';
import View from 'ol/View.js';

document.addEventListener('DOMContentLoaded', function() {
  const map = new Map({
    target: 'map',
    layers: [
      new TileLayer({
        title: 'OpenTopoMap',
        type: 'base',
        visible: true,
        source: new XYZ({
          url: 'https://{a|b|c}.tile.opentopomap.org/{z}/{x}/{y}.png'
        })
      }),
    ],
    view: new View({
      center: [0, 0],
      zoom: 2,
    }),
  });

}, false);


