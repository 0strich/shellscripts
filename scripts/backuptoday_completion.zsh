#!/usr/bin/env zsh

_backuptoday() {
    local -a opts
    opts=("saltmine" "vendor" "develop" "ilharu" "tria")
    _describe 'values' opts
}

compdef _backuptoday backuptoday
