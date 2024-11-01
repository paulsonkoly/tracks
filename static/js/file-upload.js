document.addEventListener('DOMContentLoaded', function() {
  const fileUpload = document.getElementById('file-upload');
  const fileName = document.getElementById('file-name');
  const fileError = document.getElementById('file-error');
  const dropArea = document.getElementById('drop-area');
  const submitButton = document.getElementById('file-submit');

  // Initially disable the submit button
  submitButton.disabled = true;
  submitButton.classList.add('bg-gray-400', 'cursor-not-allowed');

  // Function to enable/disable the button based on file input
  function toggleSubmitButton(isEnabled) {
    submitButton.disabled = !isEnabled;
    submitButton.classList.toggle('bg-green-600', isEnabled);
    submitButton.classList.toggle('bg-gray-400', !isEnabled);
    submitButton.classList.toggle('cursor-not-allowed', !isEnabled);
    submitButton.classList.toggle('hover:bg-green-700', isEnabled);
  }

  // File input event listener
  fileUpload.addEventListener('change', () => {
    const file = fileUpload.files[0];
    if (file && file.name.endsWith('.gpx')) {
      fileName.textContent = `Selected file: ${file.name}`;
      fileName.classList.remove('hidden');
      fileError.classList.add('hidden');
      toggleSubmitButton(true); // Enable the submit button
    } else {
      fileUpload.value = ''; // Reset file input
      fileError.textContent = 'Please upload a valid .gpx file.';
      fileError.classList.remove('hidden');
      fileName.classList.add('hidden');
      toggleSubmitButton(false); // Disable the submit button
    }
  });

  // Drag-and-drop area events
  dropArea.addEventListener('click', () => fileUpload.click());
  dropArea.addEventListener('dragover', (e) => {
    e.preventDefault();
    dropArea.classList.add('bg-green-50');
  });
  dropArea.addEventListener('dragleave', () => dropArea.classList.remove('bg-green-50'));
  dropArea.addEventListener('drop', (e) => {
    e.preventDefault();
    dropArea.classList.remove('bg-green-50');
    fileUpload.files = e.dataTransfer.files;
    fileUpload.dispatchEvent(new Event('change'));
  });
});

