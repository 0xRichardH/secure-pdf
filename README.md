## Secure PDF

![image](https://github.com/0x4richard/secure-pdf/assets/7600503/9e1fb122-7278-460c-a813-876c2acb9092)


- packaging for desktop

  https://developer.fyne.io/started/packaging

```
go install fyne.io/fyne/v2/cmd/fyne@latest
fyne package -os darwin -icon icon.png
```

https://superuser.com/questions/526920/how-to-remove-quarantine-from-file-permissions-in-os-x

```
ls -l@
xattr -r -d com.apple.quarantine secure_pdf.app
```
