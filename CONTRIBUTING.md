# Contribute to the lucky - community

Welcome to Lucky, here is a list of contribution guides for you. If you find something incorrect or missing on a page, please submit an issue or PR to fix it.


## What can you do
Encourage any action you take to make the project better. On GitHub, every improvement to a project can be made through a PR (short for pull request).


* If you find a typo, try fixing it!
* If you find a bug, try to fix it!
* If you find some redundant code, please try to delete them!
* If you find that some test cases are missing, try adding them!
* If you can enhance the functionality, don't hesitate!
* If you find that code is implicit, try adding comments to make it clear!
* If you think the code is ugly, try refactoring it!
* If you can help improve the documentation, that would be great!
* If you find the documentation is incorrect, go ahead and fix it!
  *...




## contribute
### equipment
Before contributing, you need to register a Github ID. Prepare the following environment:
* JDK 1.8
  *git

### Workflow

1. Fork to your repository
2. Clone the fork to your local repository
3. Create a new branch and work on it
4. Keep your branches in sync
5. Commit the changes (make sure the commit message is concise)
6. Push commits to the forked repository
7. Create a PR

Please follow the pull request template. Please make sure the PR has a corresponding issue.
After a PR is created, one or more reviewers are assigned to the pull request. Reviewers will review the code.
Please remove any commits that fix review feedback, typos, merges, and rebases before merging the PR. The final commit message should be concise and clear



### Submit rules
#### Submit message

Commit messages can help reviewers better understand what the purpose of submitting a PR is. It can also help speed up the code review process. Contributors are encouraged to submit messages using EXPLICIT rather than ambiguous messages. In general, we advocate the following commit message types:

* feat: new features
* fix: fix errors
* docs: change documentation
* style: changes that do not affect the meaning of the code (spaces, formatting, missing semicolons, etc.)
* refactor: Refactoring: code changes that neither fix bugs nor add functionality
* perf: Code changes to improve performance
* test: add missing tests or correct existing tests
* chore: changes to the build process or auxiliary tools and libraries (such as documentation generation) (chore: 7788)

On the other hand, we discourage contributors from submitting messages via:

* ~~fix bug Fix bug~~
* ~~update update~~
* ~~add doc Add document~~




#### Submit content

A commit represents all content changes contained in a commit. It is best for us to include the content in a single submission, which supports the reviewer's complete review without the help of any other submissions. In other words, content within a single commit can be passed through CI to avoid code clutter. In short, there are two small rules we need to keep in mind:

* Avoid very large changes in commits;
* Every submission is complete and viewable.

Whether it's commit messages or commit content, we do take code review more seriously.


### Pull Request

Note that a single PR should not be too large. If a large number of changes are required, it is best to separate the changes into several separate PRs.

### Code Review
All code should be carefully reviewed by one or more committers. Some principles:

- Readability: Important code should be well documented. Follow our coding style.
- Elegance: New functions, classes or components should be well designed.
- Testability: Important code should be fully tested (high unit test coverage).



### Sign your work
A signature is a simple line at the end of a patch's explanation that certifies that you wrote it or have the right to pass it on as an open source patch. The rules are very simple: if you can prove the following (from [developercertificate.org](http://developercertificate.org/)):

```
Developer Certificate of Origin
Version 1.1

Copyright (C) 2004, 2006 The Linux Foundation and its contributors.
660 York Street, Suite 102,
San Francisco, CA 94110 USA

Everyone is permitted to copy and verbatim copies of this distribute
license document, but changing it is not allowed.

Developer's Certificate of Origin 1.1

By making a contribution to this project, I certify that:

(a) The contribution was created in whole or in part by me and I
     have the right to submit it under the open source license
     or indicated in the file;

(b) The contribution is based upon previous work that, to the best
     of my knowledge, is covered under an appropriate open source
     license and I have the right under that license to submit that
     work with modifications, whether created in whole or in part
     by me, under the same open source license (unless I am
     permitted to submit under a different license), as indicated
     in the file; or

(c) The contribution was provided directly to me by some other
     person who certified (a), (b) or (c) and I have not modified
     it.

(d) I understand and agree that this project and the contribution
     are public and that a record of the contribution (including all
     personal information I submit with it, including my sign-off) is
     maintained indefinitely and may be redistributed consistently with
     this project or the open source license(s) involved.
```

Then you just add a line to each git commit message:

```
Signed-off-by: Joe Smith <joe.smith@email.com>
```

Use your real name (sorry, no pseudonyms or anonymous contributions).

If the user.name and user.email git configuration are set, you can use git commit -s to automatically sign commits.