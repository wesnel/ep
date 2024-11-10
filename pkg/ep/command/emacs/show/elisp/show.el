;;; show.el --- Use `auth-source' to show a password  -*- lexical-binding: t; -*-

;; Copyright (C) 2024 Wesley Nelson <wgn@wgn.dev>

;; URL: https://github.com/wesnel/ep
;; Author: Wesley Nelson <wgn@wgn.dev>
;; Maintainer: Wesley Nelson <wgn@wgn.dev>
;; Created: 09 Nov 2024
;; Keywords: convenience

;; Package-Requires: ((emacs "23.1"))

;; This file is not part of GNU Emacs.

;; This file is free software: you can redistribute it and/or modify
;; it under the terms of the GNU General Public License as published by
;; the Free Software Foundation, either version 3 of the License, or
;; (at your option) any later version.

;; This file is distributed in the hope that it will be useful,
;; but WITHOUT ANY WARRANTY; without even the implied warranty of
;; MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
;; GNU General Public License for more details.

;; You should have received a copy of the GNU General Public License
;; along with this file.  If not, see <https://www.gnu.org/licenses/>.

;;; Commentary:

;;; Code:

(require 'auth-source)

(condition-case err
    (let* ((host "{{.Host}}")
           (user "{{.User}}")
           (port (string-to-number "{{.Port}}"))
           (spec '())
           (spec (if (not (eq host ""))
                     (append '(:host host) spec)))
           (spec (if (not (eq user ""))
                     (append '(:user user) spec)))
           (spec (if (not (eq port 0))
                     (append '(:port port) spec)))
           (password
            (auth-info-password
             (nth 0
                  (apply #'auth-source-search spec)))))
      (when password
        (princ password)))
  (error
   (print err #'external-debugging-output)))

(provide 'show)

;;; show.el ends here
