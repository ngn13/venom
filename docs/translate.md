## Adding/Updating translations
Start by forking and cloning the repo. Then create a new branch and 
change directory into the server:
```bash
cd server/
```

Now make sure you have python installed:
```bash
python3 --version
```
The binary may also be named `python` and not `python3`.

Then run the `update-lang.py` script by specifying a 2 letter country
code for the language. For example:
```bash
python3 update-lang.py tr
```

This will update the translation if it already exists, and it will create 
if it does not exist. Lastly you can edit the translation, it will be located 
under the `lang/` folder.

When you are done, commit your changes and create a [pull request](https://github.com/ngn13/venom/pulls).

