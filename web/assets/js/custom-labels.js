// Custom Labels Detection JavaScript

const { createApp } = Vue;

createApp({
    data() {
        return {
            selectedFile: null,
            isDetecting: false,
            labels: [],
            previewUrl: null,
            fileInfo: null,
            error: null
        };
    },
    mounted() {
        this.setupFileUpload();
    },
    methods: {
        setupFileUpload() {
            const dropZone = this.$refs.uploadArea;
            const fileInput = this.$refs.fileInput;
            
            FileHandler.setupDragAndDrop(dropZone, fileInput, (file) => {
                this.handleFileSelect(file);
            });
        },

        async handleFileSelect(file) {
            try {
                // Validate file
                FileHandler.validateImage(file);
                
                this.selectedFile = file;
                this.fileInfo = {
                    name: file.name,
                    size: Utils.formatFileSize(file.size),
                    type: file.type
                };
                
                // Create preview
                const previewContainer = this.$refs.previewContainer;
                this.previewUrl = await FileHandler.createImagePreview(file, previewContainer);
                
                this.error = null;
                this.labels = [];
                
            } catch (error) {
                this.error = error.message;
                this.clearFile();
            }
        },

        clearFile() {
            this.selectedFile = null;
            this.previewUrl = null;
            this.fileInfo = null;
            this.labels = [];
            this.$refs.previewContainer.innerHTML = '';
            this.$refs.fileInput.value = '';
        },

        async detectLabels() {
            if (!this.selectedFile) {
                this.error = 'Please select an image first';
                return;
            }

            this.isDetecting = true;
            this.error = null;
            this.labels = [];

            try {
                const query = `
                    mutation UploadAndDetectCustomLabels($file: Upload!) {
                        uploadAndDetectCustomLabels(file: $file) {
                            labels {
                                name
                                confidence
                            }
                        }
                    }
                `;

                const files = {
                    file: this.selectedFile
                };

                const result = await GraphQL.query(query, {}, files);
                this.labels = result.uploadAndDetectCustomLabels.labels;

                if (this.labels.length === 0) {
                    this.error = 'No custom labels detected in this image.';
                }

            } catch (error) {
                console.error('Error detecting labels:', error);
                this.error = 'Failed to detect labels: ' + error.message;
            } finally {
                this.isDetecting = false;
            }
        },

        getConfidenceColor(confidence) {
            if (confidence >= 80) return '#28a745';
            if (confidence >= 60) return '#ffc107';
            return '#dc3545';
        },

        getConfidenceText(confidence) {
            if (confidence >= 80) return 'High';
            if (confidence >= 60) return 'Medium';
            return 'Low';
        },

        reset() {
            this.clearFile();
            this.error = null;
        }
    }
}).mount('#app');