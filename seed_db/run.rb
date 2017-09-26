require 'typhoeus'
require 'json'
require 'dotenv'
require 'commander/import'

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
    seed_repo(options.organization, options.repo)
  end
end

BASE_URL = 'https://api.github.com'.freeze
BASE_OPTIONS = { userpwd: "jfo84:#{ENV.fetch('GITHUB_ACCESS_TOKEN')}" }

def seed_repo(organization_id, repo_id)
  record_repo(organization_id, repo_id)
  puts "Seeding database for #{organization_id}/#{repo_id}"
  current_page = 1
  loop do
    puts "Starting page #{current_page}..."
    pulls_request = Typhoeus::Request.new("#{BASE_URL}/repos/#{organization_id}/#{repo_id}/pulls", 
      params: { page: current_page, state: 'all' },
      **BASE_OPTIONS
    )
    pulls_request.run
    pulls_hash = JSON.parse(pulls_request.response.body)

    pull_urls = pulls_hash.map { |pull_hash| pull_hash['url'] }
  
    pull_urls.each do |pull_url|
      pull = record_pull(pull_url, repo_id)
      # Since comments and reviews hang off pulls
      next if pull.is_dup?
      record_comments(pull)
      record_reviews(pull)
      # Safeguard to avoid rate limits. I don't want to deal with the empty responses
      # Limit is 5000/min
      sleep 0.5
    end
  
    break if pull_urls.length < 30
    current_page += 1
  end
  puts 'Done'
end

def record_repo(organization_id, repo_id)
  repo_request = Typhoeus::Request.new("#{BASE_URL}/repos/#{organization_id}/#{repo_id}", **BASE_OPTIONS)
  repo_request.run
  repo_hash = JSON.parse(repo_request.response.body)
  repo = Repo.new(data_hash: repo_hash)
  repo.record unless repo.is_dup?
end

def record_pull(pull_url, repo_id)
  # I wish there was a better name for this ^_^
  pull_request = Typhoeus::Request.new(pull_url, **BASE_OPTIONS)
  pull_request.run
  pull_hash = JSON.parse(pull_request.response.body)
  pull = Pull.new(data_hash: pull_hash, repo_id: repo_id)
  puts pull.id
  pull.record unless pull.is_dup?
  pull
end

def record_comments(pull)
  comments_request = Typhoeus::Request.new("#{pull.url}/comments", **BASE_OPTIONS)
  comments_request.run
  comments = JSON.parse(comments_request.response.body)

  comments.each do |comment_hash|
    comment = Comment.new(data_hash: comment_hash, pull_id: pull.id)
    comment.record
  end
end

def record_reviews(pull)
  reviews_request = Typhoeus::Request.new("#{pull.url}/reviews", **BASE_OPTIONS)
  reviews_request.run
  reviews = JSON.parse(reviews_request.response.body)

  reviews.each do |review_hash|
    review_json = review_hash.to_json
    review = Review.new(data_hash: review_hash, pull_id: pull.id)
    review.record unless review.is_dup?
  end
end
