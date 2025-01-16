build:
	pushd frontend; npm run build; popd
	cp -R frontend/dist/frontend/browser/* backend/static/
	pushd backend; go build; mv programmerjournal-backend ../; popd;
