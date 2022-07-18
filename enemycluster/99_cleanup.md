# Cleanup/Spuren verwischen
```
unset HISTFILE
echo "Nothing to see here, sorry" > ~/.bash_history
find / -exec touch {} +
```

* `sudo journalctl --vacuum-time=1s` - logs löschen

4. Timestamps von modifierten Files resetten  
  ```
  touch -d "date"
  ```
