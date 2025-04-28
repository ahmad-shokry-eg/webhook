make push:
	@read -p "Enter commit message: " message; \
	git add .; \
	git commit -m "$$message";
	GIT_SSH_COMMAND='ssh -i ~/.ssh/id_ed25519' git push origin main