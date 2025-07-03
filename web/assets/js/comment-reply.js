// Comment Reply Generator JavaScript

const { createApp } = Vue;

createApp({
    data() {
        return {
            selectedFile: null,
            originalComment: '',
            isGenerating: false,
            replies: [],
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
                
            } catch (error) {
                this.error = error.message;
                this.clearFile();
            }
        },

        clearFile() {
            this.selectedFile = null;
            this.previewUrl = null;
            this.fileInfo = null;
            this.$refs.previewContainer.innerHTML = '';
            this.$refs.fileInput.value = '';
        },

        async generateReplies() {
            if (!this.selectedFile || !this.originalComment.trim()) {
                this.error = 'Please select an image and enter a comment';
                return;
            }

            this.isGenerating = true;
            this.error = null;
            this.replies = [];

            try {
                const query = `
                    mutation GenerateCommentReplies($input: GenerateCommentRepliesInput!, $file: Upload!) {
                        generateCommentReplies(input: $input, file: $file) {
                            replies {
                                style
                                content
                            }
                        }
                    }
                `;

                const variables = {
                    input: {
                        originalComment: this.originalComment
                    }
                };

                const files = {
                    file: this.selectedFile
                };

                const result = await GraphQL.query(query, variables, files);
                this.replies = result.generateCommentReplies.replies;

                if (this.replies.length === 0) {
                    this.error = 'No replies were generated. Please try again.';
                }

            } catch (error) {
                console.error('Error generating replies:', error);
                this.error = 'Failed to generate replies: ' + error.message;
            } finally {
                this.isGenerating = false;
            }
        },

        async copyReply(content) {
            const success = await Utils.copyToClipboard(content);
            if (success) {
                Utils.showSuccess('Reply copied to clipboard!');
            } else {
                Utils.showError('Failed to copy reply');
            }
        },

        reset() {
            this.clearFile();
            this.originalComment = '';
            this.replies = [];
            this.error = null;
        },

        getStyleColor(style) {
            const colors = {
                friendly: '#28a745',
                professional: '#007bff',
                humorous: '#ffc107'
            };
            return colors[style.toLowerCase()] || '#667eea';
        },

        getStyleIcon(style) {
            const icons = {
                friendly: '😊',
                professional: '💼',
                humorous: '😄'
            };
            return icons[style.toLowerCase()] || '💬';
        }
    }
}).mount('#app');