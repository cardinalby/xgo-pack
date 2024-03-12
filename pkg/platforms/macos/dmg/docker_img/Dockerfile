FROM sporsh/create-dmg:latest

COPY create-dmg-with-app_symlink.sh .

ENV APPLICATIONS_SYMLINK=1

ENTRYPOINT ["/bin/bash", "/create-dmg-with-app_symlink.sh"]