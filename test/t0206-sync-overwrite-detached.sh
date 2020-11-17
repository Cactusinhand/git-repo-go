#!/bin/sh

test_description="sync overwrites modified files in detached head"

. ./lib/sharness.sh

# Create manifest repositories
manifest_url="file://${REPO_TEST_REPOSITORIES}/hello/manifests"

test_expect_success "setup" '
	# create .repo file as a barrier, not find .repo deeper
	touch .repo &&
	mkdir work
'

test_expect_success "git-repo sync to Maint branch" '
	(
		cd work &&
		git-repo init -u $manifest_url -b Maint &&
		git-repo sync \
			--mock-ssh-info-status 200 \
			--mock-ssh-info-response \
			"{\"host\":\"ssh.example.com\", \"port\":22, \"type\":\"agit\"}"
	)
'

test_expect_success "manifests version: 1.0" '
	(
		cd work &&
		cat >expect<<-EOF &&
		manifests: Version 1.0
		EOF
		(
			cd .repo/manifests &&
			git log -1 --pretty="manifests: %s"
		) >actual &&
		test_cmp expect actual
	)
'

test_expect_success "edit files in workdir, all projects are in detached HEAD" '
	(
		cd work &&
		test -f drivers/driver-1/VERSION &&
		echo hacked >drivers/driver-1/VERSION &&
		test -f projects/app1/VERSION &&
		echo hacked >projects/app1/VERSION &&
		test -f projects/app1/module1/VERSION &&
		echo hacked >projects/app1/module1/VERSION &&
		test -f projects/app2/VERSION &&
		echo hacked >projects/app2/VERSION
	)
'

test_expect_success "git-repo sync to master branch, do not overwrite edit files" '
	(
		cd work &&
		git-repo init -u $manifest_url -b master &&
		test_must_fail git-repo sync \
			--mock-ssh-info-status 200 \
			--mock-ssh-info-response \
			"{\"host\":\"ssh.example.com\", \"port\":22, \"type\":\"agit\"}" &&
		(
			cd drivers/driver-1 &&
			git status --porcelain &&
			cd ../../projects/app1 &&
			git status --porcelain &&
			cd ../../projects/app2 &&
			git status --porcelain &&
			cd ../../projects/app1/module1 &&
			git status --porcelain
		) >actual &&
		cat >expect <<-EOF &&
		 M VERSION
		 M VERSION
		?? module1/
		 M VERSION
		 M VERSION
		EOF
		test_cmp expect actual
	)
'

test_expect_success "clean workspace" '
	(
		cd work &&
		(
			cd drivers/driver-1 &&
			git checkout -- . &&
			cd ../../projects/app1 &&
			git checkout -- . &&
			cd ../../projects/app2 &&
			git checkout -- . &&
			cd ../../projects/app1/module1 &&
			git checkout -- .
		) &&
		git-repo init -u $manifest_url -b master &&
		git-repo sync \
			--mock-ssh-info-status 200 \
			--mock-ssh-info-response \
			"{\"host\":\"ssh.example.com\", \"port\":22, \"type\":\"agit\"}"
	)
'

test_expect_success "manifests version: 2.0" '
	(
		cd work &&
		cat >expect<<-EOF &&
		manifests: Version 2.0
		EOF
		(
			cd .repo/manifests &&
			git log -1 --pretty="manifests: %s"
		) >actual &&
		test_cmp expect actual
	)
'

test_expect_success "projects after sync" '
	(
		cd work &&
		cat >expect<<-EOF &&
		driver-1: v1.0-dev
		app-1: v2.0.0-dev
		app-1.module1: v1.0.0
		app-2: v2.0.0-dev
		EOF
		echo "driver-1: $(cat drivers/driver-1/VERSION)" >actual &&
		echo "app-1: $(cat projects/app1/VERSION)" >>actual &&
		echo "app-1.module1: $(cat projects/app1/module1/VERSION)" >>actual &&
		echo "app-2: $(cat projects/app2/VERSION)" >>actual &&
		test_cmp expect actual
	)
'

test_done
