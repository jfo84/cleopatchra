require 'json'
require 'dotenv'
require 'commander/import'

require_relative './request'
require_relative './repo'
require_relative './pull'
require_relative './comment'
require_relative './review'

Dotenv.load

program :name, 'Database Seeder'
program :version, '0.0.1'
program :description, 'Seeds the database'

command :seed do |c|
  c.syntax = 'seed [options]'
  c.description = 'Seeds the database for a given repo'
  c.option '--organization STRING', String, 'Specifies an organization for the repo'
  c.option '--repo STRING', String, 'Specifies the repo that we want to choose'
  c.action do |args, options|
    # TODO: Remove the default
    options.default :organization => 'facebook', :repo => 'react'
    repo_external_id = "#{options.organization}/#{options.repo}"
    seed_repo(repo_external_id)
  end
end

BASE_URL = 'https://api.github.com'.freeze
BASE_OPTIONS = { userpwd: "jfo84:#{ENV.fetch('GITHUB_ACCESS_TOKEN')}" }

def seed_repo(repo_external_id)
  repo = record_repo(repo_external_id)
  puts "Seeding database for #{repo_external_id}"
  current_page = 1
  loop do
    puts "Starting page #{current_page}..."
    pulls_request = Request.new("#{BASE_URL}/repos/#{repo_external_id}/pulls", 
      params: { page: current_page, state: 'all' },
      **BASE_OPTIONS
    )
    pulls_request.run
    pulls_hash = JSON.parse(pulls_request.response.body)

    pull_urls = pulls_hash.map { |pull_hash| pull_hash['url'] }
  
    pull_urls.each do |pull_url|
      pull = record_pull(pull_url, repo.id)
      # Since comments and reviews hang off pulls
      next if pull.is_dup?
      record_comments(pull)
      record_reviews(pull)
    end
  
    break if pull_urls.length < 30
    current_page += 1
  end
  puts 'Done'
end

def record_repo(repo_external_id)
  repo_request = Request.new("#{BASE_URL}/repos/#{repo_external_id}", **BASE_OPTIONS)
  repo_request.run
  repo_hash = JSON.parse(repo_request.response.body)
  repo = Repo.new(data_hash: repo_hash)
  repo.record unless repo.is_dup?
  repo
end

def record_pull(pull_url, repo_id)
  # I wish there was a better name for this ^_^
  pull_request = Request.new(pull_url, **BASE_OPTIONS)
  pull_request.run
  pull_hash = JSON.parse(pull_request.response.body)
  pull = Pull.new(data_hash: pull_hash, repo_id: repo_id)
  puts pull.id
  pull.record unless pull.is_dup?
  pull
end

def record_comments(pull)
  comments_request = Request.new(pull.comments_url, **BASE_OPTIONS)
  comments_request.run
  comments = JSON.parse(comments_request.response.body)

  comments.each do |comment_hash|
    comment = Comment.new(data_hash: comment_hash, pull_id: pull.id)
    comment.record unless comment.is_dup?
  end
end

def record_reviews(pull)
  reviews_request = Request.new("#{pull.url}/reviews", **BASE_OPTIONS)
  reviews_request.run
  reviews = JSON.parse(reviews_request.response.body)

  reviews.each do |review_hash|
    review_json = review_hash.to_json
    review = Review.new(data_hash: review_hash, pull_id: pull.id)
    review.record unless review.is_dup?
  end
end
