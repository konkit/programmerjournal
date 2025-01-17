build:
	pushd frontend; npm install; npm run build; popd
	cp -R frontend/dist/frontend/browser/* backend/static/
	pushd backend; go build -o programmerjournal; mv programmerjournal ../; popd;
