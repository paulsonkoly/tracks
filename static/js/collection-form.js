document.addEventListener('DOMContentLoaded', () => {
  const tracksContainer = document.getElementById('tracksContainer');
  const addTrackButton = document.getElementById('addTrackButton');
  let trackIndex = 0;

  // Store added track IDs for validation
  const addedTrackIds = new Set();

  // Add a new track input field
  addTrackButton.addEventListener('click', () => {
    trackIndex++;
    const trackField = document.createElement('div');
    trackField.className = 'relative track-field';
    trackField.innerHTML = `
      <label class="block text-gray-600 font-medium">Track Name</label>
      <div class="flex space-x-2 items-center">
        <input type="text" 
          class="track-input w-full px-4 py-2 border rounded-md focus:outline-none focus:ring-2 focus:ring-blue-600 focus:border-transparent"
          placeholder="Enter track name" autocomplete="off">
        <input type="hidden" name="track_ids[]" class="track-id">
        <button type="button" class="delete-track bg-red-600 text-white font-semibold px-2 py-1 rounded transition hover:bg-red-700">
          Delete
        </button>
      </div>
      <ul class="suggestions hidden absolute z-10 bg-white border rounded-md shadow-lg mt-1 max-h-48 overflow-y-auto w-full"></ul>
    `;
    tracksContainer.appendChild(trackField);
  });

  // Event delegation for dynamic input fields
  tracksContainer.addEventListener('input', async (event) => {
    if (event.target.classList.contains('track-input')) {
      const inputField = event.target;
      const suggestionsBox = inputField.closest('.track-field').querySelector('.suggestions');
      const hiddenInput = inputField.closest('.track-field').querySelector('.track-id');
      const query = inputField.value.trim();

      // If the input field is cleared, remove its ID from the addedTrackIds set
      if (!query) {
        const oldId = hiddenInput.value;
        if (oldId) {
          addedTrackIds.delete(parseInt(oldId, 10));
          hiddenInput.value = ""; // Clear the hidden input
        }
        suggestionsBox.classList.add('hidden');
        return;
      }

      // Skip requests for query lengths less than 3
      if (query.length < 3) {
        suggestionsBox.classList.add('hidden');
        return;
      }

      try {
        const response = await fetch(`/tracks?name=${query}`);
        if (response.ok) {
          const data = await response.json();
          suggestionsBox.innerHTML = '';

          // Filter out already-added track IDs
          const availableTracks = data.tracks.filter((track) => !addedTrackIds.has(track.id));

          if (availableTracks.length === 0) {
            suggestionsBox.classList.add('hidden');
            return;
          }

          availableTracks.forEach((track) => {
            const listItem = document.createElement('li');
            listItem.textContent = track.name;
            listItem.className = 'px-4 py-2 hover:bg-blue-100 cursor-pointer';
            listItem.addEventListener('click', () => {
              const oldId = hiddenInput.value;
              if (oldId) {
                addedTrackIds.delete(parseInt(oldId, 10)); // Remove the old ID
              }
              inputField.value = track.name;
              hiddenInput.value = track.id;
              addedTrackIds.add(track.id);
              suggestionsBox.classList.add('hidden');
            });
            suggestionsBox.appendChild(listItem);
          });

          suggestionsBox.classList.remove('hidden');
        }
      } catch (error) {
        console.error('Error fetching track suggestions:', error);
      }
    }
  });

  // Delete track field
  tracksContainer.addEventListener('click', (event) => {
    if (event.target.classList.contains('delete-track')) {
      const trackField = event.target.closest('.track-field');
      const hiddenInput = trackField.querySelector('.track-id');
      const trackId = hiddenInput.value;

      if (trackId) {
        addedTrackIds.delete(parseInt(trackId, 10));
      }

      tracksContainer.removeChild(trackField);
    }
  });

  // Hide suggestions when clicking outside
  document.addEventListener('click', (event) => {
    if (!event.target.closest('.track-field')) {
      document.querySelectorAll('.suggestions').forEach((box) => {
        box.classList.add('hidden');
      });
    }
  });
});

