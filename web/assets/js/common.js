// Common JavaScript utilities and functions

// Utility functions
const Utils = {
    // Show loading state
    showLoading(element, text = 'Loading...') {
        if (typeof element === 'string') {
            element = document.querySelector(element);
        }
        if (element) {
            element.innerHTML = `<span class="loading"></span> ${text}`;
            element.disabled = true;
        }
    },

    // Hide loading state
    hideLoading(element, originalText) {
        if (typeof element === 'string') {
            element = document.querySelector(element);
        }
        if (element) {
            element.innerHTML = originalText;
            element.disabled = false;
        }
    },

    // Show error message
    showError(message, container = document.body) {
        const errorDiv = document.createElement('div');
        errorDiv.className = 'error';
        errorDiv.textContent = message;
        container.insertBefore(errorDiv, container.firstChild);
        
        // Auto remove after 5 seconds
        setTimeout(() => {
            if (errorDiv.parentNode) {
                errorDiv.parentNode.removeChild(errorDiv);
            }
        }, 5000);
    },

    // Show success message
    showSuccess(message, container = document.body) {
        const successDiv = document.createElement('div');
        successDiv.className = 'success';
        successDiv.textContent = message;
        container.insertBefore(successDiv, container.firstChild);
        
        // Auto remove after 3 seconds
        setTimeout(() => {
            if (successDiv.parentNode) {
                successDiv.parentNode.removeChild(successDiv);
            }
        }, 3000);
    },

    // Format file size
    formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    },

    // Format date
    formatDate(date) {
        return new Date(date).toLocaleString();
    },

    // Debounce function
    debounce(func, wait) {
        let timeout;
        return function executedFunction(...args) {
            const later = () => {
                clearTimeout(timeout);
                func(...args);
            };
            clearTimeout(timeout);
            timeout = setTimeout(later, wait);
        };
    },

    // Copy text to clipboard
    async copyToClipboard(text) {
        try {
            await navigator.clipboard.writeText(text);
            return true;
        } catch (err) {
            console.error('Failed to copy text: ', err);
            return false;
        }
    }
};

// GraphQL API utilities
const GraphQL = {
    async query(query, variables = {}, files = {}) {
        try {
            if (Object.keys(files).length > 0) {
                // File upload with GraphQL multipart spec
                const formData = new FormData();
                
                // Prepare operations
                const operations = {
                    query: query,
                    variables: { ...variables }
                };

                // Create map for file uploads
                const map = {};
                let fileIndex = 0;
                
                // Process files and create map
                for (const [key, file] of Object.entries(files)) {
                    map[fileIndex.toString()] = [`variables.${key}`];
                    operations.variables[key] = null;
                    fileIndex++;
                }
                
                // Add operations and map to form data
                formData.append('operations', JSON.stringify(operations));
                formData.append('map', JSON.stringify(map));
                
                // Add files to form data
                fileIndex = 0;
                for (const [key, file] of Object.entries(files)) {
                    formData.append(fileIndex.toString(), file);
                    fileIndex++;
                }

                const response = await fetch('/query', {
                    method: 'POST',
                    body: formData
                });

                const result = await response.json();
                
                if (result.errors) {
                    throw new Error(result.errors[0].message);
                }
                
                return result.data;
            } else {
                // Regular GraphQL query without files
                const response = await fetch('/query', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        query: query,
                        variables: variables
                    })
                });

                const result = await response.json();
                
                if (result.errors) {
                    throw new Error(result.errors[0].message);
                }
                
                return result.data;
            }
        } catch (error) {
            console.error('GraphQL Error:', error);
            throw error;
        }
    }
};

// File handling utilities
const FileHandler = {
    // Validate image file
    validateImage(file) {
        const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];
        const maxSize = 10 * 1024 * 1024; // 10MB
        
        if (!validTypes.includes(file.type)) {
            throw new Error('Please select a valid image file (JPEG, PNG, GIF, or WebP)');
        }
        
        if (file.size > maxSize) {
            throw new Error('File size must be less than 10MB');
        }
        
        return true;
    },

    // Create image preview
    createImagePreview(file, container) {
        return new Promise((resolve, reject) => {
            const reader = new FileReader();
            
            reader.onload = function(e) {
                const img = document.createElement('img');
                img.src = e.target.result;
                img.className = 'image-preview';
                img.alt = 'Preview';
                
                // Clear container and add image
                container.innerHTML = '';
                container.appendChild(img);
                
                resolve(e.target.result);
            };
            
            reader.onerror = function() {
                reject(new Error('Failed to read file'));
            };
            
            reader.readAsDataURL(file);
        });
    },

    // Setup drag and drop
    setupDragAndDrop(dropZone, fileInput, onFileSelected) {
        // Prevent default drag behaviors
        ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, preventDefaults, false);
            document.body.addEventListener(eventName, preventDefaults, false);
        });

        // Highlight drop zone when item is dragged over it
        ['dragenter', 'dragover'].forEach(eventName => {
            dropZone.addEventListener(eventName, () => dropZone.classList.add('dragover'), false);
        });

        ['dragleave', 'drop'].forEach(eventName => {
            dropZone.addEventListener(eventName, () => dropZone.classList.remove('dragover'), false);
        });

        // Handle dropped files
        dropZone.addEventListener('drop', handleDrop, false);
        
        // Handle click to select file
        dropZone.addEventListener('click', () => fileInput.click());
        
        // Handle file input change
        fileInput.addEventListener('change', (e) => {
            if (e.target.files.length > 0) {
                onFileSelected(e.target.files[0]);
            }
        });

        function preventDefaults(e) {
            e.preventDefault();
            e.stopPropagation();
        }

        function handleDrop(e) {
            const dt = e.dataTransfer;
            const files = dt.files;
            
            if (files.length > 0) {
                onFileSelected(files[0]);
            }
        }
    }
};

// Initialize common functionality when DOM is loaded
document.addEventListener('DOMContentLoaded', function() {
    // Add smooth scrolling to all anchor links
    document.querySelectorAll('a[href^="#"]').forEach(anchor => {
        anchor.addEventListener('click', function (e) {
            e.preventDefault();
            const target = document.querySelector(this.getAttribute('href'));
            if (target) {
                target.scrollIntoView({
                    behavior: 'smooth'
                });
            }
        });
    });
    
    // Add focus management for accessibility
    document.addEventListener('keydown', function(e) {
        if (e.key === 'Tab') {
            document.body.classList.add('keyboard-navigation');
        }
    });
    
    document.addEventListener('mousedown', function() {
        document.body.classList.remove('keyboard-navigation');
    });
});

// Export utilities for use in other scripts
window.Utils = Utils;
window.GraphQL = GraphQL;
window.FileHandler = FileHandler;