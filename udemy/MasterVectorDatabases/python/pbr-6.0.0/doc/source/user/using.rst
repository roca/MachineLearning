=======
 Usage
=======

*pbr* is a *setuptools* plugin and so to use it you must use *setuptools* and
call ``setuptools.setup()``. While the normal *setuptools* facilities are
available, *pbr* makes it possible to express them through static data files.

.. _setup_py:

``setup.py``
------------

*pbr* only requires a minimal ``setup.py`` file compared to a standard
*setuptools* project. This is because most configuration is located in static
configuration files. This recommended minimal ``setup.py`` file should look
something like this::

    #!/usr/bin/env python

    from setuptools import setup

    setup(
        setup_requires=['pbr'],
        pbr=True,
    )

.. note::

   It is necessary to specify ``pbr=True`` to enabled *pbr* functionality.

.. note::

   While one can pass any arguments supported by setuptools to ``setup()``,
   any conflicting arguments supplied in ``setup.cfg`` will take precedence.

``pyproject.toml``
------------------

PBR can be configured as a PEP517 build-system in ``pyproject.toml``. This
currently continues to rely on setuptools which means you need the above
``setup.py`` file to be present. The main benefits to using a
``pyproject.toml`` file with PBR are that you can control the versions of
PBR and setuptools that are used avoiding easy_install invocation.
Your build-system block in ``pyproject.toml`` will need to look something
like this::

    [build-system]
    requires = ["pbr>=6.0.0", "setuptools>=64.0.0"]
    build-backend = "pbr.build"

Eventually PBR may grow its own direct support for PEP517 build hooks, but
until then it will continue to need setuptools and ``setup.py``.

.. _setup_cfg:

``setup.cfg``
-------------

The ``setup.cfg`` file is an INI-like file that can mostly replace the
``setup.py`` file. It is similar to the ``setup.cfg`` file found in recent
versions of `setuptools`__. A simple sample can be found in *pbr*'s own
``setup.cfg`` (it uses its own machinery to install itself):

::

    [metadata]
    name = pbr
    author = OpenStack Foundation
    author_email = openstack-discuss@lists.openstack.org
    summary = OpenStack's setup automation in a reusable form
    description_file = README.rst
    description_content_type = text/x-rst; charset=UTF-8
    home_page = https://launchpad.net/pbr
    project_urls =
        Bug Tracker = https://bugs.launchpad.net/pbr/
        Documentation = https://docs.openstack.org/pbr/
        Source Code = https://opendev.org/openstack/pbr
    license = Apache-2
    classifier =
        Development Status :: 4 - Beta
        Environment :: Console
        Environment :: OpenStack
        Intended Audience :: Developers
        Intended Audience :: Information Technology
        License :: OSI Approved :: Apache Software License
        Operating System :: OS Independent
        Programming Language :: Python
    keywords =
        setup
        distutils

    [files]
    packages =
        pbr
    data_files =
        etc/pbr = etc/*
        etc/init =
            pbr.packaging.conf
            pbr.version.conf

    [entry_points]
    console_scripts =
        pbr = pbr.cmd:main
    pbr.config.drivers =
        plain = pbr.cfg.driver:Plain

Recent versions of `setuptools`_ provide many of the same sections as *pbr*.
However, *pbr* does provide a number of additional sections:

- ``files``
- ``entry_points``
- ``backwards_compat``
- ``pbr``

In addition, there are some modifications to other sections:

- ``metadata``

For all other sections, you should refer to either the `setuptools`_
documentation or the documentation of the package that provides the section,
such as the ``extract_messages`` section provided by Babel__.

.. note::

   Comments may be used in ``setup.cfg``, however all comments should start
   with a ``#`` and may be on a single line, or in line, with at least one
   white space character immediately preceding the ``#``. Semicolons are not a
   supported comment delimiter. For instance::

       [section]
       # A comment at the start of a dedicated line
       key =
           value1 # An in line comment
           value2
           # A comment on a dedicated line
           value3

.. note::

   On Python 3 ``setup.cfg`` is explicitly read as UTF-8.  On Python 2 the
   encoding is dependent on the terminal encoding.

__ http://setuptools.readthedocs.io/en/latest/setuptools.html#configuring-setup-using-setup-cfg-files
__ http://babel.pocoo.org/en/latest/setup.html

``files``
~~~~~~~~~

The ``files`` section defines the install location of files in the package
using three fundamental keys: ``packages``, ``namespace_packages``, and
``data_files``.

``packages``
  A list of top-level packages that should be installed. The behavior of
  packages is similar to ``setuptools.find_packages`` in that it recurses the
  Python package hierarchy below the given top level and installs all of it. If
  ``packages`` is not specified, it defaults to the value of the ``name`` field
  given in the ``[metadata]`` section.

``namespace_packages``
  Similar to ``packages``, but is a list of packages that provide namespace
  packages.

``data_files``
  A list of files to be installed. The format is an indented block that
  contains key value pairs which specify target directory and source file to
  install there. More than one source file for a directory may be indicated
  with a further indented list. Source files are stripped of leading
  directories. Additionally, *pbr* supports a simple file globbing syntax for
  installing entire directory structures. For example::

      [files]
      data_files =
          etc/pbr = etc/pbr/*
          etc/neutron =
              etc/api-paste.ini
              etc/dhcp-agent.ini
          etc/init.d = neutron.init

  This will result in ``/etc/neutron`` containing ``api-paste.ini`` and
  ``dhcp-agent.ini``, both of which *pbr* will expect to find in the ``etc``
  directory in the root of the source tree. Additionally, ``neutron.init`` from
  that directory will be installed in ``/etc/init.d``. All of the files and
  directories located under ``etc/pbr`` in the source tree will be installed
  into ``/etc/pbr``.

  Note that this behavior is relative to the effective root of the environment
  into which the packages are installed, so depending on available permissions
  this could be the actual system-wide ``/etc`` directory or just a top-level
  ``etc`` subdirectory of a *virtualenv*.

``entry_points``
~~~~~~~~~~~~~~~~

The ``entry_points`` section defines entry points for generated console scripts
and Python libraries. This is actually provided by *setuptools* but is
documented here owing to its importance.

The general syntax of specifying entry points is a top level name indicating
the entry point group name, followed by one or more key value pairs naming
the entry point to be installed. For instance::

    [entry_points]
    console_scripts =
        pbr = pbr.cmd:main
    pbr.config.drivers =
        plain = pbr.cfg.driver:Plain
        fancy = pbr.cfg.driver:Fancy

Will cause a console script called *pbr* to be installed that executes the
``main`` function found in ``pbr.cmd``. Additionally, two entry points will be
installed for ``pbr.config.drivers``, one called ``plain`` which maps to the
``Plain`` class in ``pbr.cfg.driver`` and one called ``fancy`` which maps to
the ``Fancy`` class in ``pbr.cfg.driver``.

``backwards_compat``
~~~~~~~~~~~~~~~~~~~~~

.. todo:: Describe this section

.. _pbr-setup-cfg:

``pbr``
~~~~~~~

The ``pbr`` section controls *pbr*-specific options and behaviours.

``skip_git_sdist``
  If enabled, *pbr* will not generate a manifest file from *git* commits. If
  this is enabled, you may need to define your own `manifest template`__.

  This can also be configured using the ``SKIP_GIT_SDIST`` environment
  variable, as described :ref:`here <packaging-tarballs>`.

  __ https://packaging.python.org/tutorials/distributing-packages/#manifest-in

``skip_changelog``
  If enabled, *pbr* will not generated a ``ChangeLog`` file from *git* commits.

  This can also be configured using the ``SKIP_WRITE_GIT_CHANGELOG``
  environment variable, as described :ref:`here <packaging-authors-changelog>`

``skip_authors``
  If enabled, *pbr* will not generate an ``AUTHORS`` file from *git* commits.

  This can also be configured using the ``SKIP_GENERATE_AUTHORS`` environment
  variable, as described :ref:`here <packaging-authors-changelog>`

``skip_reno``
  If enabled, *pbr* will not generate a ``RELEASENOTES.txt`` file if `reno`_ is
  present and configured.

  This can also be configured using the ``SKIP_GENERATE_RENO`` environment
  variable, as described :ref:`here <packaging-releasenotes>`.

``autodoc_tree_index_modules``
  A boolean option controlling whether *pbr* should generate an index of
  modules using ``sphinx-apidoc``. By default, all files except ``setup.py``
  are included, but this can be overridden using the ``autodoc_tree_excludes``
  option.

  .. deprecated:: 4.2

      This feature has been replaced by the `sphinxcontrib-apidoc`_ extension.
      Refer to the :ref:`build_sphinx` overview for more information.

``autodoc_tree_excludes``
  A list of modules to exclude when building documentation using
  ``sphinx-apidoc``. Defaults to ``[setup.py]``. Refer to the
  `sphinx-apidoc man page`__ for more information.

  __ http://sphinx-doc.org/man/sphinx-apidoc.html

  .. deprecated:: 4.2

      This feature has been replaced by the `sphinxcontrib-apidoc`_ extension.
      Refer to the :ref:`build_sphinx` overview for more information.

``autodoc_index_modules``
  A boolean option controlling whether *pbr* should itself generates
  documentation for Python modules of the project. By default, all found Python
  modules are included; some of them can be excluded by listing them in
  ``autodoc_exclude_modules``.

  .. deprecated:: 4.2

      This feature has been replaced by the `sphinxcontrib-apidoc`_ extension.
      Refer to the :ref:`build_sphinx` overview for more information.

``autodoc_exclude_modules``
  A list of modules to exclude when building module documentation using *pbr*.
  *fnmatch* style pattern (e.g. ``myapp.tests.*``) can be used.

  .. deprecated:: 4.2

      This feature has been replaced by the `sphinxcontrib-apidoc`_ extension.
      Refer to the :ref:`build_sphinx` overview for more information.

``api_doc_dir``
  A subdirectory inside the ``build_sphinx.source_dir`` where auto-generated
  API documentation should be written, if ``autodoc_index_modules`` is set to
  True. Defaults to ``"api"``.

  .. deprecated:: 4.2

      This feature has been replaced by the `sphinxcontrib-apidoc`_ extension.
      Refer to the :ref:`build_sphinx` overview for more information.

.. note::

   When using ``autodoc_tree_excludes`` or ``autodoc_index_modules`` you may
   also need to set ``exclude_patterns`` in your Sphinx configuration file
   (generally found at ``doc/source/conf.py`` in most OpenStack projects)
   otherwise Sphinx may complain about documents that are not in a toctree.
   This is especially true if the ``[sphinx_build] warning-is-error`` option is
   set. See the `Sphinx build configuration file`__ documentation for more
   information on configuring Sphinx.

   __ http://sphinx-doc.org/config.html

.. versionchanged:: 4.2

   The ``autodoc_tree_index_modules``, ``autodoc_tree_excludes``,
   ``autodoc_index_modules``, ``autodoc_exclude_modules`` and ``api_doc_dir``
   settings are all deprecated.

.. versionchanged:: 2.0

   The ``pbr`` section used to take a ``warnerrors`` option that would enable
   the ``-W`` (Turn warnings into errors.) option when building Sphinx. This
   feature was broken in 1.10 and was removed in pbr 2.0 in favour of the
   ``[build_sphinx] warning-is-error`` provided in Sphinx 1.5+.

``metadata``
~~~~~~~~~~~~

.. todo:: Describe this section

.. _build_sphinx-setup-cfg:

``build_sphinx``
~~~~~~~~~~~~~~~~

.. versionchanged:: 3.0

   The ``build_sphinx`` plugin used to default to building both HTML and man
   page output. This is no longer the case, and you should explicitly set
   ``builders`` to ``html man`` if you wish to retain this behavior.

.. deprecated:: 4.2

   This feature has been superseded by the `sphinxcontrib-apidoc`_ (for
   generation of API documentation) and :ref:`pbr.sphinxext` (for configuration
   of versioning via package metadata) extensions. It has been removed in
   version 6.0.

Requirements
------------

Requirements files are used in place of the ``install_requires`` and
``extras_require`` attributes. Requirement files should be given one of the
below names. This order is also the order that the requirements are tried in:

* ``requirements.txt``
* ``tools/pip-requires``

Only the first file found is used to install the list of packages it contains.

.. versionchanged:: 5.0

   Previously you could specify requirements for a given major version of
   Python using requirements files with a ``-pyN`` suffix. This was deprecated
   in 4.0 and removed in 5.0 in favour of environment markers.

.. _extra-requirements:

Extra requirements
~~~~~~~~~~~~~~~~~~

Groups of optional dependencies, or `"extra" requirements`__, can be described
in your ``setup.cfg``, rather than needing to be added to ``setup.py``. An
example (which also demonstrates the use of environment markers) is shown
below.

__ https://www.python.org/dev/peps/pep-0426/#extras-optional-dependencies

Environment markers
~~~~~~~~~~~~~~~~~~~

Environment markers are `conditional dependencies`__ which can be added to the
requirements (or to a group of extra requirements) automatically, depending on
the environment the installer is running in. They can be added to requirements
in the requirements file, or to extras defined in ``setup.cfg``, but the format
is slightly different for each.

For ``requirements.txt``::

    argparse; python_version=='2.6'

This will result in the package depending on ``argparse`` only if it's being
installed into Python 2.6.

For extras specified in ``setup.cfg``, add an ``extras`` section. For instance,
to create two groups of extra requirements with additional constraints on the
environment, you can use::

    [extras]
    security =
        aleph
        bet:python_version=='3.2'
        gimel:python_version=='2.7'
    testing =
        quux:python_version=='2.7'

__ https://www.python.org/dev/peps/pep-0426/#environment-markers

Testing
-------

.. deprecated:: 4.0

As described in :doc:`/user/features`, *pbr* may override the ``test`` command
depending on the test runner used.

A typical usage would be in ``tox.ini`` such as::

  [tox]
  minversion = 2.0
  skipsdist = True
  envlist = py33,py34,py35,py26,py27,pypy,pep8,docs

  [testenv]
  usedevelop = True
  setenv =
    VIRTUAL_ENV={envdir}
    CLIENT_NAME=pbr
  deps = .
       -r{toxinidir}/test-requirements.txt
  commands =
    python setup.py test --testr-args='{posargs}'

The argument ``--coverage`` will set ``PYTHON`` to ``coverage run`` to produce
a coverage report.  ``--coverage-package-name`` can be used to modify or narrow
the packages traced.


Sphinx ``conf.py``
------------------

As described in :doc:`/user/features`, *pbr* provides a Sphinx extension to
automatically configure the version numbers for your documentation using *pbr*
metadata.

To enable this extension, you must add it to the list of extensions in
your ``conf.py`` file::

    extensions = [
        'pbr.sphinxext',
        # ... other extensions
    ]

You should also unset/remove the ``version`` and ``release`` attributes from
this file.

.. _setuptools: http://www.sphinx-doc.org/en/stable/setuptools.html
.. _sphinxcontrib-apidoc: https://pypi.org/project/sphinxcontrib-apidoc/
.. _reno: https://docs.openstack.org/reno/latest/
